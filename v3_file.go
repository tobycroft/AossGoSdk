package AossGoSdk

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	Net "github.com/tobycroft/TuuzNet"
)

// V3 开放接口（AOSSTP8 /v3/file/token/*）Go 客户端
//
// 功能与旧版文件上传临时凭证（File，/v2/file/token/*）一致：
//   - GetUploadToken           换取一次性上传 token
//   - GetUploadUrl             换取一次性 token + 完整上传地址
//   - GetUploadHashUrl         换取 hash 模式上传地址（上传后仅返回 md5）
//   - CheckMd5                  秒传前置：按 md5 查系统内是否已有文件（跨项目）
//   - PrepareUpload/PrepareUploadByMd5 命中直接返回 path，未命中返回标准上传地址
//   - GetUploadedFileUrlByHash 按 md5 查询文件完整信息（当前 Bucket 分组）
//   - Delete / Restore          软删（移入 upload/recycle/YYYYMMDD/）与恢复
//
// 鉴权改为 V3 凭证信封（见 v3_envelope.go）：请求头 X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，
// Secret 永不明文传输。存储由后台「存储管理 V3 → 存储项目」授权的 Bucket 决定：
//   - 不指定 bucket：使用该项目设定的默认 Bucket（兼容旧版调用）；
//   - 指定 bucket：必须是该项目已授权的 Bucket（按 Bucket 名称），否则签发接口返回 403。
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

// v3BucketPost 不传 bucket 时不带该字段，服务端自动回退默认 Bucket。
func v3BucketPost(bucket string) map[string]any {
	post := map[string]any{}
	if bucket != "" {
		post["bucket"] = bucket
	}
	return post
}

// V3Md5CheckData 秒传前置查询结果。Exists=true 时 Path/Url/Surl 为系统内既有文件。
type V3Md5CheckData struct {
	Exists bool         `json:"exists"`
	Path   string       `json:"path"`
	Url    string       `json:"url"`
	Surl   string       `json:"surl"`
	Md5    string       `json:"md5"`
	File   FileHashData `json:"file"`
}

// CheckMd5 秒传前置：按 md5 查系统内是否已有文件 POST /v3/file/token/md5_check。
// 未命中返回 Exists=false 且无 error，调用方据此决定是否继续上传。
func (self *V3File) CheckMd5(md5 string) (V3Md5CheckData, error) {
	md5 = strings.ToLower(strings.TrimSpace(md5))
	if md5 == "" {
		return V3Md5CheckData{}, errors.New("需要md5字段")
	}
	data, err := self.callToken("/v3/file/token/md5_check", map[string]any{"md5": md5})
	if err != nil {
		return V3Md5CheckData{}, err
	}
	var out V3Md5CheckData
	if err := decodeData(data, &out); err != nil {
		return V3Md5CheckData{}, err
	}
	return out, nil
}

// V3UploadPlan 上传准备结果：Hit=true 秒传命中，直接用 Path/Url；否则用 UploadUrl 提交文件。
type V3UploadPlan struct {
	Hit       bool
	Md5       string
	Path      string
	Url       string
	Surl      string
	Token     string
	ExpiredAt string
	UploadUrl string
}

// NeedUpload 未命中且已拿到上传地址时，调用方需要真正提交文件。
func (p V3UploadPlan) NeedUpload() bool {
	return !p.Hit && p.UploadUrl != ""
}

// PrepareUploadByMd5 已算好 md5 时的上传准备：命中秒传，未命中自动申请上传 token。
func (self *V3File) PrepareUploadByMd5(md5, ip, bucket string) (V3UploadPlan, error) {
	md5 = strings.ToLower(strings.TrimSpace(md5))
	if md5 == "" {
		return V3UploadPlan{}, errors.New("需要md5字段")
	}
	check, err := self.CheckMd5(md5)
	if err != nil {
		return V3UploadPlan{}, err
	}
	if check.Exists {
		return V3UploadPlan{
			Hit:  true,
			Md5:  md5,
			Path: check.Path,
			Url:  check.Url,
			Surl: check.Surl,
		}, nil
	}
	url, err := self.getUploadUrlWithIp(ip, bucket)
	if err != nil {
		return V3UploadPlan{}, err
	}
	return V3UploadPlan{
		Md5:       md5,
		Token:     url.Token,
		ExpiredAt: url.ExpiredAt,
		UploadUrl: url.UploadUrl,
	}, nil
}

// PrepareUpload 本地文件上传准备：先算 md5 再查系统，命中秒传，未命中返回标准上传地址。
func (self *V3File) PrepareUpload(filePath, ip, bucket string) (V3UploadPlan, error) {
	md5, err := md5LocalFile(filePath)
	if err != nil {
		return V3UploadPlan{}, err
	}
	return self.PrepareUploadByMd5(md5, ip, bucket)
}

// md5LocalFile 以流式方式计算本地文件 md5，避免大文件一次性读入内存。
func md5LocalFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetUploadToken 换取一次性上传 token POST /v3/file/token/create。
// 不指定 bucket，走服务端默认 Bucket；需指定时用 GetUploadTokenWithBucket。
func (self *V3File) GetUploadToken() (FileData, error) {
	return self.getUploadToken("")
}

// GetUploadTokenWithBucket 指定目标 Bucket（必须是该项目已授权的 Bucket 名称）。
func (self *V3File) GetUploadTokenWithBucket(bucket string) (FileData, error) {
	return self.getUploadToken(bucket)
}

func (self *V3File) getUploadToken(bucket string) (FileData, error) {
	data, err := self.callToken("/v3/file/token/create", v3BucketPost(bucket))
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
// 不指定 bucket，走服务端默认 Bucket。
func (self *V3File) GetUploadUrl() (FileUrlData, error) {
	return self.getUploadUrlWithIp("", "")
}

// GetUploadUrlWithIp 同 GetUploadUrl，额外携带上传来源 IP（写 ao_attachment_token.ip）。
func (self *V3File) GetUploadUrlWithIp(ip string) (FileUrlData, error) {
	return self.getUploadUrlWithIp(ip, "")
}

// GetUploadUrlWithBucket 指定目标 Bucket，不带来源 IP。
func (self *V3File) GetUploadUrlWithBucket(bucket string) (FileUrlData, error) {
	return self.getUploadUrlWithIp("", bucket)
}

// GetUploadUrlWithIpAndBucket 同时指定目标 Bucket 与上传来源 IP。
func (self *V3File) GetUploadUrlWithIpAndBucket(ip, bucket string) (FileUrlData, error) {
	return self.getUploadUrlWithIp(ip, bucket)
}

func (self *V3File) getUploadUrlWithIp(ip, bucket string) (FileUrlData, error) {
	post := v3BucketPost(bucket)
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
// 不指定 bucket，走服务端默认 Bucket。
func (self *V3File) GetUploadHashUrl() (FileUrlData, error) {
	return self.getUploadHashUrl("")
}

// GetUploadHashUrlWithBucket 指定目标 Bucket 的 hash 模式上传地址。
func (self *V3File) GetUploadHashUrlWithBucket(bucket string) (FileUrlData, error) {
	return self.getUploadHashUrl(bucket)
}

func (self *V3File) getUploadHashUrl(bucket string) (FileUrlData, error) {
	data, err := self.callToken("/v3/file/token/upload_url_hash", v3BucketPost(bucket))
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
