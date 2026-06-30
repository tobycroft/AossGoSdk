# 文件上传临时 Token (File Token)

## 概述

v2 版本文件上传采用临时 Token 机制，避免固定 Token 暴露给前端。

**流程：**
1. 后端通过 SDK 调用 AOSSTP8 的 `/v2/file/token/create` 获取临时 Token
2. 后端将临时 Token 返回给客户端（浏览器/App）
3. 客户端使用临时 Token 直传文件至 `/v2/file/index/upfull?token=<临时Token>`
4. 上传完成后，临时 Token 立即失效

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type FileToken struct {
    Appid     string
    Token     string
    RemoteUrl string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Appid | string | 项目 AppID |
| Token | string | 项目 Open Token（后端密钥，不可暴露给前端） |
| RemoteUrl | string | 可选，自定义远程地址，不填则使用默认地址 |

---

## GetUploadToken

获取临时上传 Token。

```go
func (self *FileToken) GetUploadToken() (FileTokenData, error)
```

**返回值 `FileTokenData`：**

```go
type FileTokenData struct {
    Token     string `json:"token"`
    ExpiredAt string `json:"expired_at"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Token | string | 临时上传 Token |
| ExpiredAt | string | 过期时间 |

---

## 使用示例

### Go 后端获取临时 Token

```go
package main

import (
    "fmt"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func getUploadToken(w http.ResponseWriter, r *http.Request) {
    ft := &AossGoSdk.FileToken{
        Appid: "your-appid",
        Token: "your-open-token",
    }

    data, err := ft.GetUploadToken()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // 将临时 token 返回给前端
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"token":"%s","expired_at":"%s"}`, data.Token, data.ExpiredAt)
}
```

### 前端使用临时 Token 上传

```javascript
// 1. 从后端获取临时 token
const resp = await fetch('/api/upload/token');
const { token } = await resp.json();

// 2. 使用临时 token 直传文件
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const uploadResp = await fetch(
    `https://upload.tuuz.cc:444/v2/file/index/upfull?token=${token}`,
    { method: 'POST', body: formData }
);

const result = await uploadResp.json();
console.log('上传结果:', result.data);
```

---

## 签名机制

SDK 内部自动生成签名：

```
sign = MD5(Appid + Token + timestamp)
```

- `timestamp` 为 Unix 时间戳
- 服务端验证签名和时间戳有效性（5 分钟内）
- 临时 Token 有效期 5 分钟，一次性使用

---

## 完整流程示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

var ft = &AossGoSdk.FileToken{
    Appid: "your-appid",
    Token: "your-open-token",
}

// 后端接口：前端请求获取临时上传 token
func uploadTokenHandler(w http.ResponseWriter, r *http.Request) {
    data, err := ft.GetUploadToken()
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "code": -1,
            "msg":  err.Error(),
        })
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 0,
        "data": map[string]string{
            "token":      data.Token,
            "expired_at": data.ExpiredAt,
            "upload_url": "https://upload.tuuz.cc:444/v2/file/index/upfull",
        },
    })
}

func main() {
    http.HandleFunc("/api/upload/token", uploadTokenHandler)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```