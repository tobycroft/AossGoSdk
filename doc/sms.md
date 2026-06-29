# 短信 (SMS)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type ASMS struct {
    Name  string
    Token string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 项目名称/标识 |
| Token | string | 项目 Token，用于签名 |

---

## Sms_send

发送单条短信。

```go
func (self *ASMS) Sms_send(phone any, quhao, text, ip any) error
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| phone | any | 手机号码 |
| quhao | string | 区号，如 "86" |
| text | string | 短信内容 |
| ip | any | 客户端 IP 地址 |

**返回值：**

| 返回值 | 说明 |
|--------|------|
| error | nil 表示发送成功，非 nil 表示发送失败 |

**使用示例：**

```go
package main

import (
    "fmt"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func main() {
    sms := &AossGoSdk.ASMS{
        Name:  "your-project-name",
        Token: "your-project-token",
    }

    err := sms.Sms_send("13800138000", "86", "您的验证码是 123456", "192.168.1.1")
    if err != nil {
        fmt.Println("短信发送失败:", err)
        return
    }
    fmt.Println("短信发送成功")
}
```

## 签名机制

SDK 内部会自动生成签名，签名算法为：

```
sign = MD5(Token + timestamp)
```

其中 `timestamp` 为 Unix 时间戳，参数通过 POST form-data 方式发送至 `/v1/sms/single/push`。