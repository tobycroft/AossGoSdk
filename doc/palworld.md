# 幻兽帕鲁服务器管理 (PalWorld)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type PalWorld struct {
    Name  string
    Token string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 项目名称/标识 |
| Token | string | 项目 Token，用于签名 |

---

## 数据类型

```go
type OnlineUser struct {
    Name      string `json:"name"`
    Playeruid string `json:"playeruid"`
    Steamid   string `json:"steamid"`
}
```

---

## 方法列表

| 方法 | 说明 |
|------|------|
| ShowPlayers | 获取在线玩家列表 |
| Kick | 踢出玩家 |
| Ban | 封禁玩家 |
| Ping | 检测服务器连通性 |

---

## ShowPlayers

获取当前在线玩家列表。

```go
func (self PalWorld) ShowPlayers() ([]OnlineUser, error)
```

**使用示例：**

```go
pal := AossGoSdk.PalWorld{
    Name:  "your-project-name",
    Token: "your-project-token",
}

players, err := pal.ShowPlayers()
if err != nil {
    fmt.Println("获取玩家列表失败:", err)
    return
}

fmt.Printf("当前在线玩家: %d 人\n", len(players))
for _, p := range players {
    fmt.Printf("  玩家: %s, UID: %s, SteamID: %s\n", p.Name, p.Playeruid, p.Steamid)
}
```

---

## Kick

踢出指定玩家。

```go
func (self PalWorld) Kick(id any) error
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | any | 玩家 ID（Playeruid 或 Steamid） |

**使用示例：**

```go
pal := AossGoSdk.PalWorld{
    Name:  "your-project-name",
    Token: "your-project-token",
}

err := pal.Kick("player_steam_id_xxx")
if err != nil {
    fmt.Println("踢出失败:", err)
    return
}
fmt.Println("玩家已被踢出")
```

---

## Ban

封禁指定玩家。

```go
func (self PalWorld) Ban(id any) error
```

**使用示例：**

```go
err := pal.Ban("player_steam_id_xxx")
if err != nil {
    fmt.Println("封禁失败:", err)
    return
}
fmt.Println("玩家已被封禁")
```

---

## Ping

检测服务器连通性。

```go
func (self PalWorld) Ping() error
```

**使用示例：**

```go
err := pal.Ping()
if err != nil {
    fmt.Println("服务器连接失败:", err)
    return
}
fmt.Println("服务器连接正常")
```

---

## 签名机制

所有方法内部自动生成签名：

```
sign = MD5(Token + timestamp)
```

---

## 完整管理示例

```go
package main

import (
    "fmt"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func main() {
    pal := AossGoSdk.PalWorld{
        Name:  "your-project-name",
        Token: "your-project-token",
    }

    // 1. 检测服务器状态
    if err := pal.Ping(); err != nil {
        fmt.Println("服务器无法连接:", err)
        return
    }
    fmt.Println("服务器连接正常")

    // 2. 查看在线玩家
    players, err := pal.ShowPlayers()
    if err != nil {
        fmt.Println("获取玩家列表失败:", err)
        return
    }

    fmt.Printf("\n=== 在线玩家 (%d 人) ===\n", len(players))
    for _, p := range players {
        fmt.Printf("- %s (UID: %s)\n", p.Name, p.Playeruid)
    }

    // 3. 踢出指定玩家
    if len(players) > 0 {
        target := players[0]
        fmt.Printf("\n踢出玩家: %s\n", target.Name)
        if err := pal.Kick(target.Steamid); err != nil {
            fmt.Println("踢出失败:", err)
        } else {
            fmt.Println("踢出成功")
        }
    }
}
```