# 文件上传临时 Token (File Token)

## 概述

v2 版本文件上传采用临时 Token 机制，避免固定 Token 暴露给前端。

**流程：**
1. 后端通过 SDK 调用 AOSSTP8 的 `/v2/file/token/create` 获取临时 Token
2. 后端将临时 Token 返回给客户端（浏览器/App）
3. 客户端使用临时 Token 直传文件至 `/v2/file/index/upfull?token=<临时Token>`
4. 上传完成后，临时 Token 立即失效

## Go SDK

### 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

### 结构体

```go
type File struct {
    Token     string
    RemoteUrl string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Token | string | OSS Token（后端密钥，不可暴露给前端） |
| RemoteUrl | string | 可选，自定义远程地址，不填则使用默认地址 |

### 方法

```go
func (self *File) GetUploadToken() (FileData, error)
```

获取临时上传 Token。

**返回值 `FileData`：**

```go
type FileData struct {
    Token     string `json:"token"`
    ExpiredAt string `json:"expired_at"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Token | string | 临时上传 Token |
| ExpiredAt | string | 过期时间 |

```go
func (self *File) GetUploadUrl() (FileUrlData, error)
```

从 AOSSTP8 获取完整的上传地址。

**返回值 `FileUrlData`：**

```go
type FileUrlData struct {
    UploadUrl string `json:"upload_url"`
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| UploadUrl | string | 完整上传地址 |

### 使用示例

```go
package main

import (
    "fmt"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func getUploadToken(w http.ResponseWriter, r *http.Request) {
    ft := &AossGoSdk.File{
        Token: "your-oss-token",
    }

    data, err := ft.GetUploadToken()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"token":"%s","expired_at":"%s"}`, data.Token, data.ExpiredAt)
}
```

### 获取上传地址

```go
ft := &AossGoSdk.File{
    Token: "your-oss-token",
}
urlData, err := ft.GetUploadUrl()
if err == nil {
    fmt.Println(urlData.UploadUrl) // https://upload.tuuz.cc:433/v2/file/index/upfull
}
```

### 自定义地址

```go
ft := &AossGoSdk.File{
    Token:     "your-oss-token",
    RemoteUrl: "https://custom.example.com:444",
}
data, err := ft.GetUploadToken()
```

---

## PHP SDK

### 包导入

```php
use Tobycroft\AossSdk\File;
```

### 构造方法

```php
public function __construct(string $token, string $remote_url = '')
```

| 参数 | 类型 | 说明 |
|------|------|------|
| token | string | OSS Token（后端密钥，不可暴露给前端） |
| remote_url | string | 可选，自定义远程地址，不填则使用默认地址 |

### 方法

```php
public function setRemoteUrl(string $remote_url): self
```

动态修改远程地址，支持链式调用。

```php
public function getUploadToken(): FileRet
```

获取临时上传 Token。

**返回值 `FileRet`：**

| 属性 | 类型 | 说明 |
|------|------|------|
| token | string | 临时上传 Token |
| expired_at | string | 过期时间 |
| error | mixed | 错误信息，成功时为 null |
| isSuccess() | bool | 是否成功 |
| getError() | mixed | 获取错误信息 |

```php
public function getUploadUrl(): FileUrlRet
```

从 AOSSTP8 获取完整的上传地址。

**返回值 `FileUrlRet`：**

| 属性 | 类型 | 说明 |
|------|------|------|
| upload_url | string | 完整上传地址 |
| error | mixed | 错误信息，成功时为 null |
| isSuccess() | bool | 是否成功 |
| getError() | mixed | 获取错误信息 |

### 使用示例

```php
$file = new \Tobycroft\AossSdk\File('your-oss-token');
$ret = $file->getUploadToken();

if ($ret->isSuccess()) {
    echo $ret->token;       // 临时 token
    echo $ret->expired_at;  // 过期时间
} else {
    echo $ret->getError();
}
```

### 获取上传地址

```php
$file = new \Tobycroft\AossSdk\File('your-oss-token');
$ret = $file->getUploadUrl();

if ($ret->isSuccess()) {
    echo $ret->upload_url; // https://upload.tuuz.cc:433/v2/file/index/upfull
}
```

### 自定义地址

```php
// 方式一：构造函数传入
$file = new \Tobycroft\AossSdk\File('your-oss-token', 'https://custom.example.com:444');

// 方式二：链式调用
$file = (new \Tobycroft\AossSdk\File('your-oss-token'))
    ->setRemoteUrl('https://custom.example.com:444');

$ret = $file->getUploadToken();
```

---

## 前端使用临时 Token 上传

```javascript
// 1. 从后端获取临时 token
const resp = await fetch('/api/upload/token');
const { token } = await resp.json();

// 2. 使用临时 token 直传文件
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const uploadResp = await fetch(
    `https://upload.tuuz.cc:433/v2/file/index/upfull?token=${token}`,
    { method: 'POST', body: formData }
);

const result = await uploadResp.json();
console.log('上传结果:', result.data);
```

---

## 签名机制

SDK 内部自动生成签名：

```
sign = MD5(Token + timestamp)
```

- `timestamp` 为 Unix 时间戳
- 服务端验证签名和时间戳有效性（5 分钟内）
- 临时 Token 有效期 5 分钟，一次性使用

---

## 完整流程示例（Go）

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

var ft = &AossGoSdk.File{
    Token: "your-oss-token",
}

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

    urlData, _ := ft.GetUploadUrl()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 0,
        "data": map[string]string{
            "token":      data.Token,
            "expired_at": data.ExpiredAt,
            "upload_url": urlData.UploadUrl,
        },
    })
}

func main() {
    http.HandleFunc("/api/upload/token", uploadTokenHandler)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```