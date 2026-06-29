# IP 范围检测 (IP)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type IP struct {
    Code  string
    Token string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Code | string | 项目 Code |
| Token | string | 项目 Token |

---

## 方法列表

| 方法 | 说明 |
|------|------|
| IpRange | 检查 IP 是否在指定国家/省份范围内 |
| IpRangeAuth | 检查 IP 范围，不在范围内则需图形验证码验证 |
| IpRangeCaptcha | 同 IpRangeAuth，返回业务状态码 |

---

## IpRange

检查 IP 是否在指定的国家/省份范围内。

```go
func (self *IP) IpRange(country any, province []any, ip any) (bool, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| country | any | 允许的国家，如 "中国" |
| province | []any | 允许的省份列表，如 ["广东", "北京"] |
| ip | any | 客户端 IP 地址 |

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| bool | bool | true 表示在范围内 |
| error | error | 错误信息 |

**使用示例：**

```go
ipChecker := &AossGoSdk.IP{
    Code:  "your-project-code",
    Token: "your-project-token",
}

inRange, err := ipChecker.IpRange("中国", []any{"广东", "北京"}, "192.168.1.1")
if err != nil {
    fmt.Println("检测失败:", err)
    return
}
if inRange {
    fmt.Println("IP 在允许范围内")
} else {
    fmt.Println("IP 不在允许范围内")
}
```

---

## IpRangeAuth

检查 IP 是否在范围内，推荐用于短信服务等场景。如果 IP 不在范围内，则需要客户端完成图形验证码验证。

```go
func (self *IP) IpRangeAuth(country any, province []any, ip any) (bool, error)
```

**使用示例：**

```go
ipChecker := &AossGoSdk.IP{
    Code:  "your-project-code",
    Token: "your-project-token",
}

ok, err := ipChecker.IpRangeAuth("中国", []any{"广东", "北京"}, "192.168.1.1")
if err != nil {
    fmt.Println("IP 检测失败:", err)
    return
}
if !ok {
    // 需要用户完成验证码
    fmt.Println("需要验证码验证")
    return
}
// IP 验证通过，继续业务逻辑
```

---

## IpRangeCaptcha

功能与 IpRangeAuth 类似，但返回业务状态码。当 code 为 -103 时表示需要验证码。

```go
func (self *IP) IpRangeCaptcha(country any, province []any, ip any) (int64, error)
```

**返回值：**

| 返回值 | 类型 | 说明 |
|--------|------|------|
| code | int64 | 0 表示通过，-103 表示需要验证码 |
| error | error | 网络或解析错误 |

**使用示例：**

```go
ipChecker := &AossGoSdk.IP{
    Code:  "your-project-code",
    Token: "your-project-token",
}

code, err := ipChecker.IpRangeCaptcha("中国", []any{"广东", "北京"}, "192.168.1.1")
if err != nil {
    fmt.Println("检测异常:", err)
    return
}

switch code {
case 0:
    fmt.Println("IP 验证通过")
case -103:
    fmt.Println("需要验证码验证")
default:
    fmt.Printf("其他错误，code: %d\n", code)
}
```

---

## 结合验证码的完整示例

```go
package main

import (
    "fmt"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func main() {
    ipChecker := &AossGoSdk.IP{
        Code:  "your-project-code",
        Token: "your-project-token",
    }

    captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

    clientIP := "192.168.1.1"

    // 检查 IP 范围
    code, err := ipChecker.IpRangeCaptcha("中国", []any{"广东", "北京"}, clientIP)
    if err != nil {
        fmt.Println("IP 检测异常:", err)
        return
    }

    if code == -103 {
        // 需要验证码
        fmt.Println("检测到异常 IP，需要验证码验证")

        // 验证用户输入的验证码
        err := captcha.Check("session_id", "user_input_code")
        if err != nil {
            fmt.Println("验证码错误:", err)
            return
        }
        fmt.Println("验证码通过")
    }

    fmt.Println("IP 验证通过，继续业务逻辑")
}
```