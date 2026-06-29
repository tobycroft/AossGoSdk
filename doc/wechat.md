# 微信小程序 (WeChat Mini Program)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 方法列表

| 方法 | 说明 |
|------|------|
| Wechat_wxa_unlimited_file | 获取小程序二维码（302 跳转方式，推荐） |
| Wechat_wxa_unlimited_raw | 获取小程序二维码（二进制流，不推荐） |
| Wechat_wxa_scene | 解析小程序 scene 参数 |
| Wechat_sns_jscode2session | 微信授权一键登录（code2session） |
| Wechat_snsAuth | 验证 access_token 和 openid 是否匹配 |
| Wechat_wxa_getuserphonenumber | 获取用户手机号 |
| Wechat_wxa_generatescheme | 生成 scheme URL |
| Wechat_ticket_signature | 微信 JS-SDK ticket 签名 |
| Wechat_message.Send | 发送微信客服消息 |

---

## Wechat_wxa_unlimited_file

获取微信小程序二维码，返回文件 URL（推荐使用，占用少）。

```go
func Wechat_wxa_unlimited_file(project, data, page string) (string, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | string | 项目 Token |
| data | string | 小程序 scene 参数值 |
| page | string | 小程序页面路径，如 "pages/index/index" |

**返回值：** 二维码图片的 URL 地址

**使用示例：**

```go
url, err := AossGoSdk.Wechat_wxa_unlimited_file("your-project-token", "invite=123", "pages/index/index")
if err != nil {
    fmt.Println("生成二维码失败:", err)
    return
}
fmt.Println("二维码地址:", url)
```

---

## Wechat_wxa_unlimited_raw

获取微信小程序二维码，返回二进制数据（不推荐，会占用服务器内存）。

```go
func Wechat_wxa_unlimited_raw(project, data, page string) ([]byte, error)
```

**使用示例：**

```go
data, err := AossGoSdk.Wechat_wxa_unlimited_raw("your-project-token", "invite=123", "pages/index/index")
if err != nil {
    fmt.Println("生成二维码失败:", err)
    return
}
// 保存到文件
os.WriteFile("qrcode.png", data, 0644)
```

---

## Wechat_wxa_scene

解析小程序 scene 参数。

```go
func Wechat_wxa_scene(project, scene string) (WechatWxaScene, error)
```

**返回值 `WechatWxaScene`：**

```go
type WechatWxaScene struct {
    Key  string
    Val  string
    Page string
    Path string
    Url  string
}
```

**使用示例：**

```go
scene, err := AossGoSdk.Wechat_wxa_scene("your-project-token", "invite=abc123")
if err != nil {
    fmt.Println("解析失败:", err)
    return
}
fmt.Printf("Key: %s, Val: %s, Page: %s\n", scene.Key, scene.Val, scene.Page)
```

---

## Wechat_sns_jscode2session

微信授权一键登录，通过前端获取的 js_code 换取 session_key 和 openid。

```go
func Wechat_sns_jscode2session(project, js_code string) (WechatSnsJscode2session, error)
```

**返回值 `WechatSnsJscode2session`：**

```go
type WechatSnsJscode2session struct {
    SessionKey string
    Unionid    string
    Openid     string
}
```

**使用示例：**

```go
session, err := AossGoSdk.Wechat_sns_jscode2session("your-project-token", "081xxxx")
if err != nil {
    fmt.Println("登录失败:", err)
    return
}
fmt.Printf("OpenID: %s, SessionKey: %s\n", session.Openid, session.SessionKey)
```

---

## Wechat_snsAuth

验证前端获取的 access_token 和 openid 是否正确匹配。

```go
func Wechat_snsAuth(project, access_token, openid interface{}) (string, error)
```

**使用示例：**

```go
result, err := AossGoSdk.Wechat_snsAuth("your-project-token", "access_token_xxx", "openid_xxx")
if err != nil {
    fmt.Println("验证失败:", err)
    return
}
fmt.Println("验证通过:", result)
```

---

## Wechat_wxa_getuserphonenumber

获取用户手机号。

```go
func Wechat_wxa_getuserphonenumber(project, code string) (WechatWxaGEtUserPhoneNumber, error)
```

**返回值 `WechatWxaGEtUserPhoneNumber`：**

```go
type WechatWxaGEtUserPhoneNumber struct {
    PhoneNumber     string      `json:"phoneNumber"`
    PurePhoneNumber string      `json:"purePhoneNumber"`
    CountryCode     string      `json:"countryCode"`
    Watermark       interface{} `json:"watermark"`
}
```

**使用示例：**

```go
phone, err := AossGoSdk.Wechat_wxa_getuserphonenumber("your-project-token", "phone_code_xxx")
if err != nil {
    fmt.Println("获取手机号失败:", err)
    return
}
fmt.Printf("手机号: %s, 区号: %s\n", phone.PurePhoneNumber, phone.CountryCode)
```

---

## Wechat_wxa_generatescheme

生成微信小程序 scheme URL。

```go
func Wechat_wxa_generatescheme(project, path, query string, is_expire bool, expire_interval int) (WechatWxaGenerateScheme, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | string | 项目 Token |
| path | string | 小程序页面路径 |
| query | string | 页面参数 |
| is_expire | bool | 是否过期 |
| expire_interval | int | 过期时间间隔（天） |

**返回值 `WechatWxaGenerateScheme`：**

```go
type WechatWxaGenerateScheme struct {
    Openlink string
}
```

**使用示例：**

```go
scheme, err := AossGoSdk.Wechat_wxa_generatescheme(
    "your-project-token",
    "pages/detail/index",
    "id=123",
    true,
    30,
)
if err != nil {
    fmt.Println("生成 scheme 失败:", err)
    return
}
fmt.Println("Scheme URL:", scheme.Openlink)
```

---

## Wechat_ticket_signature

微信 JS-SDK ticket 签名。

```go
func Wechat_ticket_signature(project, noncestr string, timestamp time.Time, url string) (string, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | string | 项目 Token |
| noncestr | string | 随机字符串 |
| timestamp | time.Time | 时间戳 |
| url | string | 当前页面 URL |

**使用示例：**

```go
import "time"

sign, err := AossGoSdk.Wechat_ticket_signature(
    "your-project-token",
    "random_nonce_str",
    time.Now(),
    "https://example.com/page",
)
if err != nil {
    fmt.Println("签名失败:", err)
    return
}
fmt.Println("签名:", sign)
```

---

## Wechat_message（客服消息）

发送微信客服文本消息。

### 结构体

```go
type Wechat_message struct {
    // 私有字段，通过链式方法设置
}
```

### 方法链

```go
msg := &AossGoSdk.Wechat_message{}
msg.Set_message_text("消息内容").
    Set_openid("user_openid").
    Set_token_not_set_for_auto("project-token").
    Send()
```

**使用示例：**

```go
msg := &AossGoSdk.Wechat_message{}

err := msg.Set_message_text("您好，欢迎关注").
    Set_openid("oXXXX_user_openid").
    Set_token_not_set_for_auto("your-project-token").
    Send()

if err != nil {
    fmt.Println("发送失败:", err)
    return
}
fmt.Println("消息发送成功")
```