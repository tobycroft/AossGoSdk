package AossGoSdk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

// V3 开放接口（AOSSTP8 /v3/excel/*）Go 客户端
//
// 功能与 V1 的 token 方式 Excel 一致（上传解析 / md5 解析 / 生成下载 /
// 生成文件 URL），鉴权改为 V3 凭证信封（见 v3_envelope.go）：
// 请求头 X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，Secret 永不明文传输。
// 项目凭证由 AOSSTP8 后台「项目管理 v3」签发。
//
// 响应统一为 {code,data,echo}，code=0 时 data 为业务数据；
// CreateExcelDownloadDirectly 成功时返回 xlsx 二进制原文。

// V3Excel V3 Excel 客户端。零值可直接用，BaseUrl 为空时回退 SDK 公用地址 conf.go:baseUrl。
type V3Excel struct {
	Secret  string // 服务端签发的 64 位十六进制凭证，同时作为信封密钥
	BaseUrl string // V3 开放接口域名根；留空则复用 SDK 公用地址
}

// NewV3Excel 创建 V3 Excel 客户端。
// secret 由 AOSSTP8 后台「项目管理 v3」签发，请通过配置注入，勿硬编码。
// endpoint 留空时复用本 SDK 的公用地址（包变量 baseUrl），不在此处重复维护地址。
func NewV3Excel(secret, endpoint string) *V3Excel {
	if endpoint == "" {
		endpoint = baseUrl
	}
	return &V3Excel{Secret: secret, BaseUrl: endpoint}
}

// v3ExcelRet V3 Excel 统一响应包，Data 为 any（行集合 / URL 字符串）。
type v3ExcelRet struct {
	Code int64  `json:"code"`
	Data any    `json:"data"`
	Echo string `json:"echo"`
}

// v3BaseUrl 返回请求域名根。
func (self *V3Excel) v3BaseUrl() string {
	if self.BaseUrl == "" {
		return baseUrl
	}
	return self.BaseUrl
}

// callForm 发起一次 V3 表单请求，返回上游 data（code!=0 时返回 error）。
func (self *V3Excel) callForm(path string, post map[string]any) (any, error) {
	rets := aossPostFormAny(self.v3BaseUrl()+path, nil, post, V3AuthHeaders(self.Secret), nil)
	var resp v3ExcelRet
	if err := rets.RetJson(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Echo)
	}
	return resp.Data, nil
}

// SendExcel 上传表格文件并解析（首行作 key）POST /v3/excel/index/dp。
// 返回行集合：[]any{map[列名]单元格值}。
func (self *V3Excel) SendExcel(realPath, mimeType, fileName string) (any, error) {
	body, contentType, err := buildMultipartFile("file", realPath, mimeType, fileName)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, self.v3BaseUrl()+"/v3/excel/index/dp", body)
	if err != nil {
		return nil, err
	}
	for k, v := range V3AuthHeaders(self.Secret) {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)

	cli := aossHTTPClient(120 * time.Second)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ret v3ExcelRet
	if err := sonic.Unmarshal(raw, &ret); err != nil {
		return nil, errors.New("V3 Excel 响应格式异常: " + string(raw[:min(len(raw), 128)]))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Echo)
	}
	return ret.Data, nil
}

// SendMd5 按已上传附件 md5 解析（带服务端缓存）POST /v3/excel/search/md5。
func (self *V3Excel) SendMd5(md5 string) (any, error) {
	return self.callForm("/v3/excel/search/md5", map[string]any{"md5": md5})
}

// CreateExcelDownloadDirectly 二维数组生成 xlsx，返回二进制原文 POST /v3/excel/index/create。
// 成功时响应为 xlsx 字节流（可直接落盘/输出）；失败时响应为 JSON 错误包。
func (self *V3Excel) CreateExcelDownloadDirectly(data any) ([]byte, error) {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		return nil, err
	}
	rets := aossPostFormAny(self.v3BaseUrl()+"/v3/excel/index/create", nil,
		map[string]any{"data": payload}, V3AuthHeaders(self.Secret), nil)
	raw, err := rets.RetBytes()
	if err != nil {
		return nil, err
	}
	// xlsx 为 ZIP 二进制（"PK" 开头）；JSON 开头即视为服务端错误包
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var ret v3ExcelRet
		if err := sonic.Unmarshal(trimmed, &ret); err == nil && ret.Code != 0 {
			return nil, errors.New(ret.Echo)
		}
		return nil, errors.New("V3 Excel 返回异常数据")
	}
	return raw, nil
}

// CreateExcelFileUrl 二维数组生成 xlsx 落盘并返回可下载 URL POST /v3/excel/index/create_file。
func (self *V3Excel) CreateExcelFileUrl(data any) (string, error) {
	payload, err := sonic.MarshalString(data)
	if err != nil {
		return "", err
	}
	ret, err := self.callForm("/v3/excel/index/create_file", map[string]any{"data": payload})
	if err != nil {
		return "", err
	}
	url, ok := ret.(string)
	if !ok {
		return "", errors.New("V3 Excel 返回的 URL 格式异常")
	}
	return url, nil
}

// buildMultipartFile 以 multipart/form-data 构建单文件上传请求体。
// 返回请求体、Content-Type（含 boundary）。
func buildMultipartFile(fieldName, realPath, mimeType, fileName string) (*bytes.Buffer, string, error) {
	f, err := os.Open(realPath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	// mimeType 为空时回退 application/octet-stream；
	// CreateFormFile 不支持自定义 Content-Type，需用 CreatePart 手写 part 头
	partMime := mimeType
	if partMime == "" {
		partMime = "application/octet-stream"
	}
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(fileName)))
	hdr.Set("Content-Type", partMime)
	fw, err := w.CreatePart(hdr)
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, "", err
	}
	if err = w.Close(); err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

// escapeQuotes 转义 multipart Content-Disposition 中的引号
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
