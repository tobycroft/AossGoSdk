package AossGoSdk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/tobycroft/Calc"
)

// SNI 域名直连（域名不变，连接落到指定 IP）
//
// 内网环境访问 AOSSTP8 时，域名解析可能仍然指向公网，导致请求绕行外网。
// 开启 SNI 后，请求的 URL、Host 头、TLS ServerName 与证书校验仍使用域名，
// 只有 TCP 连接的目标地址被替换为配置的 IP，因此 HTTPS 证书校验不受影响，
// 也不需要放宽 InsecureSkipVerify。
//
// 使用示例：
//
//	AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "10.0.0.8")
//	// 域名仍为 upload.tuuz.cc，实际连接到 10.0.0.8:433
//
//	// 内网端口与外网端口不一致时，可额外指定连接端口
//	AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "10.0.0.8", 8443)
//	// 域名仍为 upload.tuuz.cc:433，实际连接到 10.0.0.8:8443
//
// 传空 IP 可移除该域名的配置：AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "")
//
// aossSniTarget 保存域名的连接目标；port 为空表示沿用 URL 中的端口。
type aossSniTarget struct {
	ip   string
	port string
}

var (
	aossSniMu sync.RWMutex
	aossSni   = map[string]aossSniTarget{}

	aossDialer = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// aossTransport 全局共享，DialContext 在建立连接时读取 SNI 映射。
	aossTransport = &http.Transport{
		DialContext:         aossDialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
)

// Aoss_conf_set_sni 为指定域名配置连接 IP（可传内网 IP 或 127.0.0.1）。
// domain 可带 scheme/端口，SDK 只取主机名；ip 传空串表示移除该域名的配置。
// port 为可选的实际连接端口：内网端口与外网不一样时传入，留空沿用 URL 中的端口。
func Aoss_conf_set_sni(domain, ip string, port ...int) {
	domain = aossNormalizeDomain(domain)
	if domain == "" {
		return
	}
	target := aossParseTarget(ip, port...)
	aossSniMu.Lock()
	if target.ip == "" {
		delete(aossSni, domain)
	} else {
		aossSni[domain] = target
	}
	aossSniMu.Unlock()
	// 让已缓存的旧连接失效，下一次请求立即使用新的目标地址。
	aossTransport.CloseIdleConnections()
}

// Aoss_conf_set_sni_map 批量设置域名到 IP 的映射，例如：
//
//	AossGoSdk.Aoss_conf_set_sni_map(map[string]string{
//	    "upload.tuuz.cc":   "10.0.0.8",
//	    "aoss.example.com": "127.0.0.1",
//	    "port.example.com": "10.0.0.9:8443", // 值支持 IP:端口
//	})
func Aoss_conf_set_sni_map(m map[string]string) {
	for domain, ip := range m {
		Aoss_conf_set_sni(domain, ip)
	}
}

// Aoss_conf_clear_sni 清除全部 SNI 配置，恢复默认 DNS 解析。
func Aoss_conf_clear_sni() {
	aossSniMu.Lock()
	aossSni = map[string]aossSniTarget{}
	aossSniMu.Unlock()
	aossTransport.CloseIdleConnections()
}

// aossSniLookup 查询域名对应的连接目标。
func aossSniLookup(host string) (aossSniTarget, bool) {
	host = aossNormalizeDomain(host)
	if host == "" {
		return aossSniTarget{}, false
	}
	aossSniMu.RLock()
	target, ok := aossSni[host]
	aossSniMu.RUnlock()
	return target, ok
}

// aossParseTarget 解析连接目标：优先使用显式 port，其次识别 "ip:port" 写法。
func aossParseTarget(ip string, port ...int) aossSniTarget {
	ip = strings.TrimSpace(ip)
	target := aossSniTarget{}
	if len(port) > 0 && port[0] > 0 {
		target.port = strconv.Itoa(port[0])
	} else if h, p, err := net.SplitHostPort(ip); err == nil {
		// 兼容 "10.0.0.8:8443" 与 "[::1]:8443"，裸 IPv6 会因多冒号解析失败而走下面的 IP 分支。
		ip = h
		target.port = p
	}
	target.ip = aossNormalizeIP(ip)
	return target
}

// aossNormalizeDomain 提取并规范化主机名：去 scheme/端口、去尾点、转小写。
func aossNormalizeDomain(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "://") {
		if u, err := url.Parse(s); err == nil && u.Hostname() != "" {
			s = u.Hostname()
		}
	} else if h, _, err := net.SplitHostPort(s); err == nil {
		s = h
	}
	s = strings.Trim(s, "[]")
	return strings.ToLower(strings.TrimSuffix(s, "."))
}

// aossNormalizeIP 规范化 IP 字面量（兼容 [::1] 写法）；非 IP 字符串原样返回。
func aossNormalizeIP(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = s[1 : len(s)-1]
	}
	return s
}

// aossDialContext 在真正建连时把域名替换为配置的 IP，TLS 握手仍使用原域名做 SNI。
func aossDialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		if target, ok := aossSniLookup(host); ok && target.ip != "" {
			if target.port != "" {
				port = target.port
			}
			addr = net.JoinHostPort(target.ip, port)
		}
	}
	return aossDialer.DialContext(ctx, network, addr)
}

// aossHTTPClient 返回共享 SNI 传输层的 HTTP 客户端。
func aossHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Transport: aossTransport, Timeout: timeout}
}

// aossResponse 保存一次请求的响应体或错误，方法名与原 TuuzNet 调用保持一致。
type aossResponse struct {
	body []byte
	err  error
}

// RetJson 把响应体解析为 JSON。
func (r *aossResponse) RetJson(v any) error {
	if r.err != nil {
		return r.err
	}
	return sonic.Unmarshal(r.body, v)
}

// RetBytes 返回原始响应体。
func (r *aossResponse) RetBytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.body, nil
}

// RetString 返回字符串响应体。
func (r *aossResponse) RetString() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return string(r.body), nil
}

// aossPostFormAny 以 multipart/form-data 发起 POST 请求。
// queries 追加到 URL 查询串，headers/cookies 可为 nil。
func aossPostFormAny(rawurl string, queries map[string]interface{}, post map[string]any, headers map[string]string, cookies map[string]string) *aossResponse {
	u, err := url.Parse(rawurl)
	if err != nil {
		return &aossResponse{err: err}
	}
	if len(queries) > 0 {
		q := u.Query()
		for k, v := range queries {
			q.Set(k, Calc.Any2String(v))
		}
		u.RawQuery = q.Encode()
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for k, v := range post {
		if err := writer.WriteField(k, Calc.Any2String(v)); err != nil {
			return &aossResponse{err: err}
		}
	}
	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return &aossResponse{err: err}
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), &body)
	if err != nil {
		return &aossResponse{err: err}
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	resp, err := aossHTTPClient(30 * time.Second).Do(req)
	if err != nil {
		return &aossResponse{err: err}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return &aossResponse{err: fmt.Errorf("读取响应失败: %w", err)}
	}
	return &aossResponse{body: raw}
}
