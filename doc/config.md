# 配置

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 默认地址

SDK 默认使用以下地址，可通过 `Wechat_conf_set` 修改：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `baseUrl` | `https://upload.tuuz.cc:433` | 基础地址 |
| `baseUrls` | `https://upload.tuuz.cc:433` | HTTPS 基础地址 |
| `cdnUrl` | `http://aoss.familyeducation.org.cn` | CDN HTTP 地址 |
| `cdnUrls` | `https://aoss.familyeducation.org.cn` | CDN HTTPS 地址 |

## Wechat_conf_set

设置全局基础地址。

```go
func Wechat_conf_set(BaseUrl, BaseUrls, CDNUrl, CDNUrls string)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| BaseUrl | string | HTTP 基础地址 |
| BaseUrls | string | HTTPS 基础地址 |
| CDNUrl | string | CDN HTTP 地址 |
| CDNUrls | string | CDN HTTPS 地址 |

**使用示例：**

```go
package main

import AossGoSdk "github.com/tobycroft/AossGoSdk"

func main() {
    AossGoSdk.Wechat_conf_set(
        "https://your-server:433",
        "https://your-server:433",
        "http://your-cdn.com",
        "https://your-cdn.com",
    )
}
```

## SNI 域名直连（内网 IP）

内网环境访问 AOSSTP8 时，域名解析可能仍指向公网。通过 SNI 配置可以保持域名不变，只把 TCP 连接指向内网 IP：

- 请求 URL、Host 头、TLS SNI 与证书校验仍使用原域名；
- 只有实际连接地址被替换为配置的 IP，不会放宽证书校验；
- 域名可带 scheme/端口，SDK 只取主机名。

```go
// 初始化阶段调用即可，对全部客户端（V1/V2/V3）生效
AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "10.0.0.8")

// 内网端口与外网端口不一致时，额外指定实际连接端口
// URL 仍为 upload.tuuz.cc:433，实际连接到 10.0.0.8:8443
AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "10.0.0.8", 8443)

// 批量设置（值支持 "IP" 或 "IP:端口"）
AossGoSdk.Aoss_conf_set_sni_map(map[string]string{
    "upload.tuuz.cc":   "10.0.0.8",
    "aoss.example.com": "127.0.0.1",
    "port.example.com": "10.0.0.9:8443",
})

// 移除单个域名（IP 传空串）
AossGoSdk.Aoss_conf_set_sni("upload.tuuz.cc", "")

// 清除全部 SNI 配置
AossGoSdk.Aoss_conf_clear_sni()
```

| 函数 | 说明 |
|------|------|
| `Aoss_conf_set_sni(domain, ip string, port ...int)` | 为域名设置连接 IP，可选指定实际连接端口；`ip` 为空时移除该域名 |
| `Aoss_conf_set_sni_map(m map[string]string)` | 批量设置，值支持 `IP` 或 `IP:端口` |
| `Aoss_conf_clear_sni()` | 清除全部配置，恢复默认 DNS 解析 |