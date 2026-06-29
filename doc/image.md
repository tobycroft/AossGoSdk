# 图片合成 (Canvas)

## 包导入

```go
import AossGoSdk "github.com/tobycroft/AossGoSdk"
```

## 结构体

```go
type Canvas struct {
    // 私有字段，通过链式方法操作
}
```

---

## 位置常量

| 常量 | 值 | 位置 |
|------|------|------|
| Canvas_Posistion_TopLeft | `"lt"` | 左上 |
| Canvas_Posistion_TopCenter | `"mt"` | 中上 |
| Canvas_Posistion_TopRight | `"rt"` | 右上 |
| Canvas_Posistion_CenterLeft | `"lm"` | 左中 |
| Canvas_Posistion_CenterCenter | `"mm"` | 居中 |
| Canvas_Posistion_CenterRight | `"rm"` | 右中 |
| Canvas_Posistion_BottomLeft | `"lb"` | 左下 |
| Canvas_Posistion_BottomCenter | `"mb"` | 中下 |
| Canvas_Posistion_BottomRight | `"rb"` | 右下 |

---

## Canvas 数据类型

```go
type Canvas_Type_Text struct {
    Type     string `json:"type"`
    Text     string `json:"text,omitempty"`
    Position string `json:"position,omitempty"`
    X        int64  `json:"x"`
    Y        int64  `json:"y"`
}

type Canvas_Type_Image struct {
    Type string `json:"type"`
    URL  string `json:"url,omitempty"`
    X    int64  `json:"x"`
    Y    int64  `json:"y"`
}
```

---

## 方法列表

| 方法 | 说明 |
|------|------|
| AddText | 添加文字图层 |
| AddImage | 添加图片图层 |
| Get_Url | 生成并获取合成图片 URL |

---

## AddText

向画布添加文字图层。

```go
func (self *Canvas) AddText(Text string, Canvas_position string, X int64, Y int64) *Canvas
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| Text | string | 文字内容 |
| Canvas_position | string | 文字位置，使用位置常量 |
| X | int64 | X 坐标 |
| Y | int64 | Y 坐标 |

**返回值：** 返回 `*Canvas` 自身，支持链式调用。

---

## AddImage

向画布添加图片图层。

```go
func (self *Canvas) AddImage(Url string, X int64, Y int64) *Canvas
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| Url | string | 图片 URL 地址 |
| X | int64 | X 坐标 |
| Y | int64 | Y 坐标 |

**返回值：** 返回 `*Canvas` 自身，支持链式调用。

---

## Get_Url

生成合成图片并返回 URL。

```go
func (self *Canvas) Get_Url(project interface{}, width int64, height int64, background_color string) (string, error)
```

**参数说明：**

| 参数 | 类型 | 说明 |
|------|------|------|
| project | interface{} | 项目 Token |
| width | int64 | 画布宽度（像素） |
| height | int64 | 画布高度（像素） |
| background_color | string | 背景颜色，如 "#FFFFFF" |

**返回值：** 合成后的图片 URL

---

## 使用示例

### 基础示例：生成带文字的图片

```go
canvas := &AossGoSdk.Canvas{}

url, err := canvas.
    AddText("Hello World", AossGoSdk.Canvas_Posistion_CenterCenter, 375, 300).
    Get_Url("your-project-token", 750, 600, "#FFFFFF")

if err != nil {
    fmt.Println("生成图片失败:", err)
    return
}
fmt.Println("图片地址:", url)
```

### 高级示例：生成海报

```go
canvas := &AossGoSdk.Canvas{}

url, err := canvas.
    // 背景图片
    AddImage("https://example.com/background.jpg", 0, 0).
    // 头像
    AddImage("https://example.com/avatar.jpg", 50, 50).
    // 用户名
    AddText("张三", AossGoSdk.Canvas_Posistion_TopLeft, 160, 60).
    // 标题
    AddText("邀请你加入课程", AossGoSdk.Canvas_Posistion_CenterCenter, 375, 300).
    // 底部文字
    AddText("扫码查看详情", AossGoSdk.Canvas_Posistion_BottomCenter, 375, 550).
    Get_Url("your-project-token", 750, 1334, "#FFFFFF")

if err != nil {
    fmt.Println("生成海报失败:", err)
    return
}
fmt.Println("海报地址:", url)
```

### 在 HTTP 接口中使用

```go
func generatePoster(w http.ResponseWriter, r *http.Request) {
    canvas := &AossGoSdk.Canvas{}

    url, err := canvas.
        AddText("专属海报", AossGoSdk.Canvas_Posistion_TopCenter, 375, 40).
        AddImage("https://example.com/qrcode.png", 275, 250).
        Get_Url("your-project-token", 750, 1334, "#FFE4B5")

    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    fmt.Fprintf(w, `{"url": "%s"}`, url)
}
```