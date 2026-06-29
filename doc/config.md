# 配置

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 默认地址

SDK 默认使用以下地址，可通过 `Wechat_conf_set` 修改：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `baseUrl` | `http://upload.tuuz.cc:81` | HTTP 基础地址 |
| `baseUrls` | `https://upload.tuuz.cc:444` | HTTPS 基础地址 |
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
        "http://your-server:81",
        "https://your-server:444",
        "http://your-cdn.com",
        "https://your-cdn.com",
    )
}
```