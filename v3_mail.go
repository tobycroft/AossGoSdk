package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
	Net "github.com/tobycroft/TuuzNet"
)

// V3 开放接口（AOSSTP8 /v3/mail/send/*）Go 客户端
//
// 服务端提供三种内容类型：
//   - SendCaptcha   验证码：主题/正文由后台模板统一维护，调用方只提供验证码（<=16 位）
//   - SendMarketing 营销：模板统一维护，调用方只替换 {变量名} 关键内容（如链接）
//   - SendCustom    自定义：主题与正文完全由调用方提供
//
// 鉴权采用 V3 凭证信封（见 v3_envelope.go），请求头 X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，
// Secret 永不明文传输。项目须在 AOSSTP8 后台授权至少一台邮件服务器；
// 使用模板的类型还须先在后台「邮件管理 → 邮件模板」配置并启用模板。
//
// 响应统一为 {code,data,echo}，code=0 时返回 data（含 message_id / type），
// 否则返回 echo 作为 error。

// V3Mail V3 邮件客户端。零值可直接用，BaseUrl 为空时回退 SDK 公用地址 conf.go:baseUrl。
type V3Mail struct {
	Secret  string // 服务端签发的 64 位十六进制凭证，同时作为信封密钥
	BaseUrl string // V3 开放接口域名根；留空则复用 SDK 公用地址
}

// NewV3Mail 创建 V3 邮件客户端。endpoint 留空时复用 SDK 公用地址。
func NewV3Mail(secret, endpoint string) *V3Mail {
	if endpoint == "" {
		endpoint = baseUrl
	}
	return &V3Mail{Secret: secret, BaseUrl: endpoint}
}

// v3MailRet V3 邮件统一响应包。
type v3MailRet struct {
	Code int64  `json:"code"`
	Data any    `json:"data"`
	Echo string `json:"echo"`
}

// v3BaseUrl 返回请求域名根。
func (self *V3Mail) v3BaseUrl() string {
	if self.BaseUrl == "" {
		return baseUrl
	}
	return self.BaseUrl
}

// call 发起一次 V3 表单请求，返回上游 data（code!=0 时返回 error）。
func (self *V3Mail) call(path string, post map[string]any) (any, error) {
	rets := new(Net.Post).PostFormDataAny(self.v3BaseUrl()+path, nil, post, V3AuthHeaders(self.Secret), nil)
	var resp v3MailRet
	if err := rets.RetJson(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Echo)
	}
	return resp.Data, nil
}

// SendCaptcha 验证码邮件 POST /v3/mail/send/captcha。
//
// template 为后台模板标识（模板须含 {code}），code 为 1~16 位验证码。
// options 支持 cc / alt_body / from_addr / from_name / server_id。
func (self *V3Mail) SendCaptcha(to, template, code string, options map[string]any) (any, error) {
	post := map[string]any{"to": to, "template": template, "code": code}
	copyOptions(post, options, "cc", "alt_body", "from_addr", "from_name", "server_id")
	return self.call("/v3/mail/send/captcha", post)
}

// SendMarketing 营销邮件 POST /v3/mail/send/marketing。
//
// vars 的键对应后台模板中的 {变量名}，最多 50 项、单值最长 1000 字符，只接受标量。
// 有占位符未提供值时会返回 error 且不发信。options 同 SendCaptcha。
func (self *V3Mail) SendMarketing(to, template string, vars map[string]any, options map[string]any) (any, error) {
	encoded, err := sonic.MarshalString(vars)
	if err != nil {
		return nil, err
	}
	post := map[string]any{"to": to, "template": template, "vars": encoded}
	copyOptions(post, options, "cc", "alt_body", "from_addr", "from_name", "server_id")
	return self.call("/v3/mail/send/marketing", post)
}

// SendCustom 自定义邮件 POST /v3/mail/send/custom。
//
// options 支持 cc / alt_body / is_html / from_addr / from_name / server_id。
func (self *V3Mail) SendCustom(to, subject, body string, options map[string]any) (any, error) {
	post := map[string]any{"to": to, "subject": subject, "body": body}
	copyOptions(post, options, "cc", "alt_body", "is_html", "from_addr", "from_name", "server_id")
	return self.call("/v3/mail/send/custom", post)
}
