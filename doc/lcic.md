# 低代码互动课堂 (LCIC)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type Lcic struct {
    Name  string
    Token string
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 项目名称/标识 |
| Token | string | 项目 Token，用于签名 |

---

## 方法列表

| 方法 | 说明 |
|------|------|
| Lcic_CreateUser | 创建用户（自动注册） |
| Lcic_RoomCreate | 创建房间 |
| Lcic_RoomModify | 修改房间信息 |
| Lcic_RoomDelete | 删除房间 |
| Lcic_LinkUrl | 获取进入房间链接 |

---

## Lcic_CreateUser

创建/注册 LCIC 用户，自动返回用户 ID 和 Token。

```go
func (self *Lcic) Lcic_CreateUser(Name, OriginId, Avatar string) (LcicStructCreateUser, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| Name | string | 用户在直播间的显示名称 |
| OriginId | string | 用户在你系统中的唯一标识符 |
| Avatar | string | 用户头像 URL 地址 |

**返回值 `LcicStructCreateUser`：**

```go
type LcicStructCreateUser struct {
    UserId string `json:"UserId"`
    Token  string `json:"Token"`
}
```

**使用示例：**

```go
lcic := &AossGoSdk.Lcic{
    Name:  "your-project-name",
    Token: "your-project-token",
}

user, err := lcic.Lcic_CreateUser(
    "张三",
    "student_001",
    "https://example.com/avatar.jpg",
)
if err != nil {
    fmt.Println("创建用户失败:", err)
    return
}
fmt.Printf("UserID: %s, Token: %s\n", user.UserId, user.Token)
```

---

## Lcic_RoomCreate

创建课堂房间。

```go
func (self *Lcic) Lcic_RoomCreate(TeacherId, StartTime, EndTime, Name any) (LcicStructCreateRoom, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| TeacherId | any | 老师 ID（LCIC 系统中的 UserId） |
| StartTime | any | 开始时间（Unix 时间戳） |
| EndTime | any | 结束时间（Unix 时间戳） |
| Name | any | 房间名称 |

**返回值 `LcicStructCreateRoom`：**

```go
type LcicStructCreateRoom struct {
    RoomId int `json:"RoomId"`
}
```

**使用示例：**

```go
import "time"

lcic := &AossGoSdk.Lcic{
    Name:  "your-project-name",
    Token: "your-project-token",
}

now := time.Now()
room, err := lcic.Lcic_RoomCreate(
    "teacher_user_id",
    now.Unix(),
    now.Add(2*time.Hour).Unix(),
    "数学课",
)
if err != nil {
    fmt.Println("创建房间失败:", err)
    return
}
fmt.Println("RoomID:", room.RoomId)
```

---

## Lcic_RoomModify

修改房间信息。

```go
func (self *Lcic) Lcic_RoomModify(RoomId, TeacherId, StartTime, EndTime, Name any) (bool, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| RoomId | any | 房间 ID |
| TeacherId | any | 老师 ID |
| StartTime | any | 开始时间 |
| EndTime | any | 结束时间 |
| Name | any | 房间名称 |

**使用示例：**

```go
ok, err := lcic.Lcic_RoomModify(
    12345,
    "teacher_user_id",
    time.Now().Unix(),
    time.Now().Add(3*time.Hour).Unix(),
    "数学课（调整时间）",
)
if err != nil {
    fmt.Println("修改房间失败:", err)
    return
}
fmt.Println("修改成功:", ok)
```

---

## Lcic_RoomDelete

删除房间。

```go
func (self *Lcic) Lcic_RoomDelete(RoomId interface{}) (bool, error)
```

**使用示例：**

```go
ok, err := lcic.Lcic_RoomDelete(12345)
if err != nil {
    fmt.Println("删除房间失败:", err)
    return
}
fmt.Println("删除成功:", ok)
```

---

## Lcic_LinkUrl

获取学生进入房间的链接。

```go
func (self *Lcic) Lcic_LinkUrl(OriginId, TeacherId interface{}) (LcicStructLinkUrl, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| OriginId | interface{} | 学生的 OriginId（你系统中的标识） |
| TeacherId | interface{} | 老师的 ID |

**返回值 `LcicStructLinkUrl`：**

```go
type LcicStructLinkUrl struct {
    Web string `json:"web"`
    Pc  string `json:"pc"`
}
```

**使用示例：**

```go
url, err := lcic.Lcic_LinkUrl("student_001", "teacher_user_id")
if err != nil {
    fmt.Println("获取链接失败:", err)
    return
}
fmt.Println("Web 端链接:", url.Web)
fmt.Println("PC 端链接:", url.Pc)
```

---

## 签名机制

所有 LCIC 方法内部自动生成签名：

```
sign = MD5(Token + timestamp)
```

---

## 完整课堂流程示例

```go
package main

import (
    "fmt"
    "time"

    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

func main() {
    lcic := &AossGoSdk.Lcic{
        Name:  "your-project-name",
        Token: "your-project-token",
    }

    // 1. 创建老师
    teacher, err := lcic.Lcic_CreateUser("王老师", "teacher_001", "https://example.com/teacher.jpg")
    if err != nil {
        fmt.Println("创建老师失败:", err)
        return
    }
    fmt.Println("老师 UserID:", teacher.UserId)

    // 2. 创建学生
    student, err := lcic.Lcic_CreateUser("小明", "student_001", "https://example.com/student.jpg")
    if err != nil {
        fmt.Println("创建学生失败:", err)
        return
    }
    fmt.Println("学生 UserID:", student.UserId)

    // 3. 创建课堂
    now := time.Now()
    room, err := lcic.Lcic_RoomCreate(
        teacher.UserId,
        now.Unix(),
        now.Add(2*time.Hour).Unix(),
        "三年级数学课",
    )
    if err != nil {
        fmt.Println("创建课堂失败:", err)
        return
    }
    fmt.Println("RoomID:", room.RoomId)

    // 4. 获取学生上课链接
    url, err := lcic.Lcic_LinkUrl("student_001", teacher.UserId)
    if err != nil {
        fmt.Println("获取链接失败:", err)
        return
    }
    fmt.Println("学生上课链接:", url.Web)
}
```