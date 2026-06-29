# 微信公众号 (Official Account)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 方法列表

| 方法 | 说明 |
|------|------|
| Wechat_offi_get_user_list | 获取已关注用户的 openid 列表 |
| Wechat_offi_get_user_info | 获取单个用户详细信息 |
| Wechat_offi_openidUrl | 获取微信网页授权 URL |
| Wechat_offi_openid_from_code | 通过 code 换取 openid |
| Wechat_template_send | 发送模板消息 |

---

## Wechat_offi_get_user_list

获取已关注用户的 openid 列表。

```go
func Wechat_offi_get_user_list(project string) ([]string, error)
```

> **注意：** 仅返回 openid 列表，无法区分用户身份。

**使用示例：**

```go
openids, err := AossGoSdk.Wechat_offi_get_user_list("your-project-token")
if err != nil {
    fmt.Println("获取用户列表失败:", err)
    return
}
for _, openid := range openids {
    fmt.Println("OpenID:", openid)
}
```

---

## Wechat_offi_get_user_info

获取单个用户的详细信息（可能获取不到）。

```go
func Wechat_offi_get_user_info(project, openid string) (WechatUserInfo, error)
```

**返回值 `WechatUserInfo`：**

```go
type WechatUserInfo struct {
    subscribe      int64
    openid         string
    nickname       string
    sex            int64
    headimgurl     string
    subscribe_time int64
}
```

**使用示例：**

```go
info, err := AossGoSdk.Wechat_offi_get_user_info("your-project-token", "oXXXX_user_openid")
if err != nil {
    fmt.Println("获取用户信息失败:", err)
    return
}
fmt.Printf("昵称: %s, 头像: %s\n", info.nickname, info.headimgurl)
```

---

## Wechat_offi_openidUrl

获取微信网页授权 URL，redirect_uri 无需手动 urlencode，SDK 自动处理。

```go
func Wechat_offi_openidUrl(project, redirect_uri, response_type, scope, state string, show_in_qrcode bool) (string, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | string | 项目 Token |
| redirect_uri | string | 回调地址（无需 urlencode） |
| response_type | string | 返回类型，通常为 "code" |
| scope | string | 授权范围，如 "snsapi_base" 或 "snsapi_userinfo" |
| state | string | 自定义状态参数 |
| show_in_qrcode | bool | 是否生成二维码显示 |

**使用示例：**

```go
url, err := AossGoSdk.Wechat_offi_openidUrl(
    "your-project-token",
    "https://example.com/callback",
    "code",
    "snsapi_userinfo",
    "state123",
    false,
)
if err != nil {
    fmt.Println("获取授权 URL 失败:", err)
    return
}
fmt.Println("授权 URL:", url)
// 将用户重定向到此 URL
```

---

## Wechat_offi_openid_from_code

通过微信授权回调返回的 code 换取 openid。

```go
func Wechat_offi_openid_from_code(project, code any) (string, error)
```

**使用示例：**

```go
// 在回调处理中
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")

    openid, err := AossGoSdk.Wechat_offi_openid_from_code("your-project-token", code)
    if err != nil {
        fmt.Println("获取 openid 失败:", err)
        return
    }
    fmt.Println("OpenID:", openid)
}
```

---

## Wechat_template_send

发送模板消息。

```go
func Wechat_template_send(project, openid, template_id, url interface{}, data map[string]Wechat_template_data_struct) (string, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | interface{} | 项目 Token |
| openid | interface{} | 用户 openid |
| template_id | interface{} | 模板 ID |
| url | interface{} | 跳转 URL |
| data | map[string]Wechat_template_data_struct | 模板数据 |

**`Wechat_template_data_struct`：**

```go
type Wechat_template_data_struct struct {
    Value string `json:"value"`
    Color string `json:"color"`
}
```

**使用示例：**

```go
data := map[string]AossGoSdk.Wechat_template_data_struct{
    "first": {
        Value: "您有一条新的消息",
        Color: "#173177",
    },
    "keyword1": {
        Value: "订单通知",
        Color: "#173177",
    },
    "keyword2": {
        Value: "2025-01-01 12:00:00",
        Color: "#173177",
    },
    "remark": {
        Value: "点击查看详情",
        Color: "#173177",
    },
}

result, err := AossGoSdk.Wechat_template_send(
    "your-project-token",
    "oXXXX_user_openid",
    "template_id_xxx",
    "https://example.com/detail",
    data,
)
if err != nil {
    fmt.Println("发送模板消息失败:", err)
    return
}
fmt.Println("发送成功:", result)
```

## 完整授权流程示例

```go
package main

import (
    "fmt"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

const projectToken = "your-project-token"

// 步骤1：引导用户授权
func Login(w http.ResponseWriter, r *http.Request) {
    url, err := AossGoSdk.Wechat_offi_openidUrl(
        projectToken,
        "https://example.com/callback",
        "code",
        "snsapi_userinfo",
        "random_state",
        false,
    )
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    http.Redirect(w, r, url, 302)
}

// 步骤2：处理回调，获取 openid
func Callback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")

    openid, err := AossGoSdk.Wechat_offi_openid_from_code(projectToken, code)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // 步骤3：获取用户信息
    info, err := AossGoSdk.Wechat_offi_get_user_info(projectToken, openid)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    fmt.Fprintf(w, "欢迎 %s", info.nickname)
}
```