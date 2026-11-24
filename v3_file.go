package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
	Net "github.com/tobycroft/TuuzNet"
)

// V3 开放接口（AOSSTP8 /v3/file/token/*）Go 客户端
//
// 功能与旧版文件上传临时凭证（File，/v2/file/token/*）一致：
//   - GetUploadToken           换取一次性上传 token
//   - GetUploadUrl             换取一次性 token + 完整上传地址
//   - GetUploadHashUrl         换取 hash 模式上传地址（上传后仅返回 md5）
//   - GetUploadedFileUrlByHash 按 md5 查询文件完整信息
//   - Delete / Restore          软删（移入 upload/recycle/YYYYMMDD/）与恢复
//
// 鉴权改为 V3 凭证信封（见 v3_envelope.go）：请求头 X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，
// Secret 永不明文传输。项目须在 AOSSTP8 后台「项目管理 v3」绑定一个 ao_oss 存储项目
// （oss_token），否则签发接口返回未配置。
//
// 前端直传阶段沿用一次性 token（upload_url），与 V2 行为一致——直传请求无法携带
// 自定义鉴权头，也不能把 secret 交给前端；信封鉴权发生在签发 token 的环节。
//
// 响应统一为 {code,data,echo}，code=0 时返回 data，否则返回 echo 作为 error。
// 响应结构与旧版 File 客户端复用（FileData / FileUrlData / FileHashData）。

// V3File V3 文件上传客户端。零值可直接用，BaseUrl 为空时回退 SDK 公用地址 conf.go:baseUrl。
type V3File struct {
	Secret  string // 服务端签发的 64 位十六进制凭证，同时作为信封密钥
	BaseUrl string // V3 开放接口域名根；留空则复用 SDK 公用地址
}

// NewV3File 创建 V3 文件上传客户端。endpoint 留空时复用 SDK 公用地址。
func NewV3File(secret, endpoint string) *V3File {
	if endpoint == "" {
		endpoint = baseUrl
	}
	return &V3File{Secret: secret, BaseUrl: endpoint}
}

// v3FileBaseUrl 返回请求域名根。
func (self *V3File) v3FileBaseUrl() string {
	if self.BaseUrl == "" {
		return baseUrl
	}
	return self.BaseUrl
}

// v3FileRet V3 file 统一响应包，Data 为 json.RawMessage 便于按方法解析。
type v3FileRet struct {
	Code int64  `json:"code"`
	Data any    `json:"data"`
	Echo string `json:"echo"`
}

// callToken 发起一次 V3 file 表单请求，返回上游 data（code!=0 时返回 error）。
func (self *V3File) callToken(path string, post map[string]any) (any, error) {
	rets := new(Net.Post).PostFormDataAny(self.v3FileBaseUrl()+path, nil, post, V3AuthHeaders(self.Secret), nil)
	var resp v3FileRet
	if err := rets.RetJson(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Echo)
	}
	return resp.Data, nil
}

// decodeData 把 callToken 返回的 data 再序列化后解析进目标结构。
func decodeData(data any, out any) error {
	raw, err := sonic.Marshal(data)
	if err != nil {
		return err
	}
	return sonic.Unmarshal(raw, out)
}

// GetUploadToken 换取一次性上传 token POST /v3/file/token/create。
func (self *V3File) GetUploadToken() (FileData, error) {
	data, err := self.callToken("/v3/file/token/create", map[string]any{})
	if err != nil {
		return FileData{}, err
	}
	var out FileData
	if err := decodeData(data, &out); err != nil {
		return FileData{}, err
	}
	return out, nil
}

// GetUploadUrl 换取一次性 token + 完整上传地址 POST /v3/file/token/upload_url。
// 可选 ip：以 map 形式写入 post（通过 SetIp 风格的直接参数）。
func (self *V3File) GetUploadUrl() (FileUrlData, error) {
	return self.getUploadUrlWithIp("")
}

// GetUploadUrlWithIp 同 GetUploadUrl，额外携带上传来源 IP（写 ao_attachment_token.ip）。
func (self *V3File) GetUploadUrlWithIp(ip string) (FileUrlData, error) {
	return self.getUploadUrlWithIp(ip)
}

func (self *V3File) getUploadUrlWithIp(ip string) (FileUrlData, error) {
	post := map[string]any{}
	if ip != "" {
		post["ip"] = ip
	}
	data, err := self.callToken("/v3/file/token/upload_url", post)
	if err != nil {
		return FileUrlData{}, err
	}
	var out FileUrlData
	if err := decodeData(data, &out); err != nil {
		return FileUrlData{}, err
	}
	return out, nil
}

// GetUploadHashUrl 换取 hash 模式上传地址 POST /v3/file/token/upload_url_hash。
func (self *V3File) GetUploadHashUrl() (FileUrlData, error) {
	data, err := self.callToken("/v3/file/token/upload_url_hash", map[string]any{})
	if err != nil {
		return FileUrlData{}, err
	}
	var out FileUrlData
	if err := decodeData(data, &out); err != nil {
		return FileUrlData{}, err
	}
	return out, nil
}

// V3FileActionData 删除/恢复接口返回的数据。
type V3FileActionData struct {
	MD5      string   `json:"md5"`
	Deleted  int      `json:"deleted"`
	Warnings []string `json:"warnings"`
}

// Delete 软删文件，服务端将文件移入 upload/recycle/YYYYMMDD/。
func (self *V3File) Delete(hash string) (V3FileActionData, error) {
	data, err := self.callToken("/v3/file/token/delete", map[string]any{"hash": hash})
	if err != nil {
		return V3FileActionData{}, err
	}
	var out V3FileActionData
	if err := decodeData(data, &out); err != nil {
		return V3FileActionData{}, err
	}
	return out, nil
}

// Restore 从回收站恢复文件到原路径。
func (self *V3File) Restore(hash string) (V3FileActionData, error) {
	data, err := self.callToken("/v3/file/token/restore", map[string]any{"hash": hash})
	if err != nil {
		return V3FileActionData{}, err
	}
	var out V3FileActionData
	if err := decodeData(data, &out); err != nil {
		return V3FileActionData{}, err
	}
	return out, nil
}

// GetUploadedFileUrlByHash 按 md5/hash 查询文件完整信息 POST /v3/file/token/hash_query。
func (self *V3File) GetUploadedFileUrlByHash(hash string) (FileHashData, error) {
	data, err := self.callToken("/v3/file/token/hash_query", map[string]any{"hash": hash})
	if err != nil {
		return FileHashData{}, err
	}
	var out FileHashData
	if err := decodeData(data, &out); err != nil {
		return FileHashData{}, err
	}
	return out, nil
}
