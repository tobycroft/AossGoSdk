# 验证码 (Captcha)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type Captcha struct {
    Token any
}
```

## 方法列表

| 方法 | 说明 |
|------|------|
| Check | 验证验证码是否正确 |
| CheckInTime | 验证验证码，并限制有效时间 |
| CheckWithCode | 验证验证码，返回业务状态码 |
| Math | 生成数学算式验证码图片 |
| Number | 生成数字验证码图片 |
| Chinese | 生成中文验证码图片 |
| Text | 生成文本验证码图片 |
| GifText | 生成动态 GIF 文本验证码（每秒一帧） |
| GifFast | 生成动态 GIF 文本验证码（0.5秒一帧） |
| GifNumber | 生成动态 GIF 数字验证码（每秒一帧） |
| GifNumberFast | 生成动态 GIF 数字验证码（0.5秒一帧） |

---

## Check

验证验证码是否正确，返回 nil 表示验证通过。

```go
func (self *Captcha) Check(ident, code any) error
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| ident | any | 验证码标识，通常使用 session id 或用户唯一标识 |
| code | any | 用户输入的验证码 |

**返回值：**

| 返回值 | 说明 |
|--------|------|
| error | nil 表示验证通过，非 nil 表示验证失败 |

**使用示例：**

```go
captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

err := captcha.Check("user_session_id", "1234")
if err != nil {
    // 验证失败
    fmt.Println("验证码错误:", err)
    return
}
// 验证通过
```

---

## CheckInTime

验证验证码，并在指定秒数内限制有效时间。

```go
func (self *Captcha) CheckInTime(ident, code any, validtime_in_second int) error
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| ident | any | 验证码标识 |
| code | any | 用户输入的验证码 |
| validtime_in_second | int | 有效期（秒），如 300 表示 5 分钟内有效 |

**使用示例：**

```go
captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

err := captcha.CheckInTime("user_session_id", "1234", 300)
if err != nil {
    fmt.Println("验证码无效或已过期:", err)
    return
}
```

---

## CheckWithCode

验证验证码，返回业务状态码和错误信息。验证码不通过时 code 为 -103。

```go
func (self *Captcha) CheckWithCode(ident, code any) (int64, error)
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| code | int64 | 业务状态码，0 表示成功，-103 表示验证码错误 |
| error | error | 网络错误或 JSON 解析错误 |

**使用示例：**

```go
captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

code, err := captcha.CheckWithCode("user_session_id", "1234")
if err != nil {
    fmt.Println("请求异常:", err)
    return
}
if code == -103 {
    fmt.Println("验证码错误")
    return
}
```

---

## 生成验证码图片

以下四个方法均返回 `image.Image` 类型的验证码图片，可直接用于 HTTP 响应或保存为文件。

### Math（数学算式）

```go
func (self *Captcha) Math(ident any) (img image.Image, err error)
```

### Number（数字）

```go
func (self *Captcha) Number(ident any) (img image.Image, err error)
```

### Chinese（中文）

```go
func (self *Captcha) Chinese(ident any) (img image.Image, err error)
```

### Text（文本）

```go
func (self *Captcha) Text(ident any) (img image.Image, err error)
```

**使用示例：**

```go
package main

import (
    "image/png"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func captchaHandler(w http.ResponseWriter, r *http.Request) {
    captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

    img, err := captcha.Math("user_session_id")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "image/png")
    png.Encode(w, img)
}
```

## 完整使用流程示例

```go
package main

import (
    "fmt"
    "image/png"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

var captcha = &AossGoSdk.Captcha{Token: "your-project-token"}

// 生成验证码
func GetCaptcha(w http.ResponseWriter, r *http.Request) {
    sessionId := "user-session-123"
    img, err := captcha.Number(sessionId)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("Content-Type", "image/png")
    png.Encode(w, img)
}

// 验证验证码
func VerifyCaptcha(w http.ResponseWriter, r *http.Request) {
    sessionId := r.FormValue("session_id")
    code := r.FormValue("code")

    err := captcha.Check(sessionId, code)
    if err != nil {
        fmt.Fprintf(w, "验证失败: %s", err)
        return
    }
    fmt.Fprintf(w, "验证通过")
}
```

---

## 动态 GIF 验证码

以下四个方法返回 `[]byte` 类型的 GIF 二进制数据，适合直接写入 HTTP 响应或保存为 `.gif` 文件。

### GifText（字母+数字，1秒一帧）

```go
func (self *Captcha) GifText(ident any) (gif []byte, err error)
```

### GifFast（字母+数字，0.5秒一帧）

```go
func (self *Captcha) GifFast(ident any) (gif []byte, err error)
```

### GifNumber（纯数字，1秒一帧）

```go
func (self *Captcha) GifNumber(ident any) (gif []byte, err error)
```

### GifNumberFast（纯数字，0.5秒一帧）

```go
func (self *Captcha) GifNumberFast(ident any) (gif []byte, err error)
```

**使用示例：**

```go
package main

import (
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func gifCaptchaHandler(w http.ResponseWriter, r *http.Request) {
    captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

    gif, err := captcha.GifNumberFast("user_session_id")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "image/gif")
    w.Write(gif)
}
```