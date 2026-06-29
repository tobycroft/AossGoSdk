# 直播 (Live Streaming)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type Live struct {
    Token string
    Code  string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Token | string | 项目 Token |
| Code | string | 项目 Code，用于签名 |

---

## 方法列表

| 方法 | 说明 |
|------|------|
| CreateAll | 创建房间并返回推流码和播放地址（全部信息） |
| CreateRoom | 创建房间并返回推流码 |
| GetPlayUrl | 获取播放地址 |

---

## CreateAll

创建直播间并返回全部信息（推流码 + 播放地址）。

```go
func (self *Live) CreateAll(title any) (LiveStructCreateAll, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| title | any | 直播间地址标识（英文数字，8 位内） |

**返回值 `LiveStructCreateAll`：**

```go
type LiveStructCreateAll struct {
    Rtmp        string `json:"rtmp"`
    Domain      string `json:"domain"`
    PlayDomain  string `json:"play_domain"`
    ObsServer   string `json:"obs_server"`
    StreamCode  string `json:"stream_code"`
    Webrtc      string `json:"webrtc"`
    Srt         string `json:"srt"`
    RtmpOverSrt string `json:"rtmp_over_srt"`
    PlayFlv     string `json:"play_flv"`
    PlayHls     string `json:"play_hls"`
    PlayRtmp    string `json:"play_rtmp"`
}
```

**使用示例：**

```go
live := &AossGoSdk.Live{
    Token: "your-project-token",
    Code:  "your-project-code",
}

all, err := live.CreateAll("room001")
if err != nil {
    fmt.Println("创建直播间失败:", err)
    return
}

fmt.Println("推流地址 (RTMP):", all.Rtmp)
fmt.Println("OBS 服务器:", all.ObsServer)
fmt.Println("推流码:", all.StreamCode)
fmt.Println("FLV 播放:", all.PlayFlv)
fmt.Println("HLS 播放:", all.PlayHls)
```

---

## CreateRoom

创建直播间并返回推流码信息（不含播放地址）。

```go
func (self *Live) CreateRoom(title any) (LiveStructCreateRoom, error)
```

**返回值 `LiveStructCreateRoom`：**

```go
type LiveStructCreateRoom struct {
    Rtmp        string `json:"rtmp"`
    Domain      string `json:"domain"`
    ObsServer   string `json:"obs_server"`
    StreamCode  string `json:"stream_code"`
    Webrtc      string `json:"webrtc"`
    Srt         string `json:"srt"`
    RtmpOverSrt string `json:"rtmp_over_srt"`
}
```

**使用示例：**

```go
live := &AossGoSdk.Live{
    Token: "your-project-token",
    Code:  "your-project-code",
}

room, err := live.CreateRoom("room002")
if err != nil {
    fmt.Println("创建房间失败:", err)
    return
}

fmt.Println("推流地址:", room.Rtmp)
fmt.Println("推流码:", room.StreamCode)
```

---

## GetPlayUrl

获取直播间播放地址。

```go
func (self *Live) GetPlayUrl(title any) (LiveStructPlayUrl, error)
```

**返回值 `LiveStructPlayUrl`：**

```go
type LiveStructPlayUrl struct {
    PlayDomain string `json:"play_domain"`
    PlayFlv    string `json:"play_flv"`
    PlayHls    string `json:"play_hls"`
    PlayRtmp   string `json:"play_rtmp"`
}
```

**使用示例：**

```go
live := &AossGoSdk.Live{
    Token: "your-project-token",
    Code:  "your-project-code",
}

url, err := live.GetPlayUrl("room001")
if err != nil {
    fmt.Println("获取播放地址失败:", err)
    return
}

fmt.Println("FLV 播放:", url.PlayFlv)
fmt.Println("HLS 播放:", url.PlayHls)
```

---

## 签名机制

SDK 内部自动生成签名：

```
sign = MD5(Code + timestamp)
```

## 完整直播流程示例

```go
package main

import (
    "fmt"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func main() {
    live := &AossGoSdk.Live{
        Token: "your-project-token",
        Code:  "your-project-code",
    }

    // 1. 创建直播间（含全部信息）
    all, err := live.CreateAll("demo001")
    if err != nil {
        fmt.Println("创建失败:", err)
        return
    }

    fmt.Println("=== 推流信息 ===")
    fmt.Println("OBS 推流地址:", all.ObsServer)
    fmt.Println("推流码:", all.StreamCode)
    fmt.Println("RTMP 推流:", all.Rtmp)

    fmt.Println("\n=== 播放信息 ===")
    fmt.Println("FLV:", all.PlayFlv)
    fmt.Println("HLS:", all.PlayHls)
    fmt.Println("RTMP:", all.PlayRtmp)
}
```