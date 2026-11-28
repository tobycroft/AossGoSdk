package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
)

// V3 开放接口（AOSSTP8 /v3/ai/*）Go 客户端
//
// 鉴权采用 V3 凭证信封（见 v3_envelope.go），一律通过请求头传递
// X-Aoss-Ts / X-Aoss-Kid / X-Aoss-Auth，不使用 V1/V2 的 ?token= 明文方式。
// 虽然服务端兼容表单/查询串传参，本客户端只用新版请求头方式。
//
// 结构化字段（messages/tools/input/thinking 等）按服务端约定以
// JSON 字符串提交；响应统一为 {code,data,echo}，code=0 时 data 为上游原样响应。

// V3Ai V3 AI 客户端。零值可直接用，BaseUrl 为空时回退 SDK 公用地址 conf.go:baseUrl。
type V3Ai struct {
	Secret  string // 服务端签发的 64 位十六进制凭证，同时作为信封密钥
	BaseUrl string // V3 开放接口域名根；留空则复用 SDK 公用地址
}

// NewV3Ai 创建 V3 AI 客户端。
// secret 由 AOSSTP8 后台「agnes」页签发，请通过配置注入，勿硬编码。
// endpoint 留空时复用本 SDK 的公用地址（包变量 baseUrl），不在此处重复维护地址。
func NewV3Ai(secret, endpoint string) *V3Ai {
	if endpoint == "" {
		endpoint = baseUrl
	}
	return &V3Ai{Secret: secret, BaseUrl: endpoint}
}

// v3Ret V3 统一响应包。
// Data 解析为 any（上游原样结构），由调用方按需要断言具体类型。
type v3Ret struct {
	Code int64  `json:"code"`
	Data any    `json:"data"`
	Echo string `json:"echo"`
}

// v3Headers 生成一次请求的 V3 鉴权请求头，复用公共 V3AuthHeaders。
func (self *V3Ai) v3Headers() map[string]string {
	return V3AuthHeaders(self.Secret)
}

// call 发起一次 V3 表单请求，返回上游 data（code!=0 时返回 error）。
func (self *V3Ai) call(path string, post map[string]any) (any, error) {
	base := self.BaseUrl
	if base == "" {
		base = baseUrl
	}
	rets := aossPostFormAny(base+path, nil, post, self.v3Headers(), nil)
	var resp v3Ret
	if err := rets.RetJson(&resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Echo)
	}
	return resp.Data, nil
}

// jsonString 将结构化字段编码为 JSON 字符串；已是字符串则原样返回。
func jsonString(v any) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return sonic.MarshalString(v)
}

// ChatCompletions 对话补全 POST /v3/ai/chat/completions。
//
// messages 可为 []map[string]any 等结构，options 支持 model/temperature/
// max_tokens/top_p/stream/tools/tool_choice/chat_template_kwargs（推理开关）。
// 推理（OpenAI 格式）：options["chat_template_kwargs"]=map[string]any{"enable_thinking":true}
func (self *V3Ai) ChatCompletions(messages any, options map[string]any) (any, error) {
	msgs, err := jsonString(messages)
	if err != nil {
		return nil, err
	}
	post := map[string]any{"messages": msgs}
	copyOptions(post, options,
		"model", "temperature", "max_tokens", "top_p", "stream", "tool_choice")
	for _, k := range []string{"tools", "chat_template_kwargs"} {
		if v, ok := options[k]; ok {
			s, err := jsonString(v)
			if err != nil {
				return nil, err
			}
			post[k] = s
		}
	}
	return self.call("/v3/ai/chat/completions", post)
}

// Responses Responses API POST /v3/ai/responses/index。
// input 可为字符串或结构化消息；推理开关走 chat_template_kwargs（OpenAI 格式）。
func (self *V3Ai) Responses(input any, options map[string]any) (any, error) {
	in, err := jsonString(input)
	if err != nil {
		return nil, err
	}
	post := map[string]any{"input": in}
	copyOptions(post, options,
		"model", "instructions", "temperature", "max_output_tokens", "top_p", "stream")
	if v, ok := options["chat_template_kwargs"]; ok {
		s, err := jsonString(v)
		if err != nil {
			return nil, err
		}
		post["chat_template_kwargs"] = s
	}
	return self.call("/v3/ai/responses/index", post)
}

// Messages Messages API（Anthropic-compatible）POST /v3/ai/messages/index。
// 推理开关按 Anthropic 格式走 thinking={"type":"enabled","budget_tokens":2048}。
func (self *V3Ai) Messages(messages any, options map[string]any) (any, error) {
	msgs, err := jsonString(messages)
	if err != nil {
		return nil, err
	}
	post := map[string]any{"messages": msgs}
	copyOptions(post, options,
		"model", "system", "temperature", "max_tokens", "top_p", "stream")
	for _, k := range []string{"thinking", "chat_template_kwargs"} {
		if v, ok := options[k]; ok {
			s, err := jsonString(v)
			if err != nil {
				return nil, err
			}
			post[k] = s
		}
	}
	return self.call("/v3/ai/messages/index", post)
}

// copyOptions 透传简单标量字段，仅复制 options 中实际存在的键。
func copyOptions(dst, src map[string]any, keys ...string) {
	for _, k := range keys {
		if v, ok := src[k]; ok {
			dst[k] = v
		}
	}
}
