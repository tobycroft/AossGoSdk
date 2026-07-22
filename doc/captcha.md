# 验证码 (Captcha)

## 概述

AOSS 提供多种验证码方案：

| 类型 | 前端交互 | 生成方法 | 验证方法 | 适用场景 |
|------|---------|---------|---------|---------|
| 文本验证码（数学/数字/中文/字母） | 用户填写文字 | `Math()` `Number()` `Chinese()` `Text()` | `Check()` `CheckInTime()` | 登录、注册、评论 |
| 动态 GIF 验证码 | 用户填写文字 | `GifText()` `GifFast()` `GifNumber()` 等 | `Check()` `CheckInTime()` | 防爬虫增强 |
| 滑动拼图验证码 | 用户拖动滑块 | `Slide()` | `SlideCheck()` | 登录、支付、重要操作 |
| 点击验证码 | 用户依次点击图标 | `Click()` | `ClickCheck()` | 高安全场景 |

---

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

---

## 一、文本验证码

### 方法列表

| 方法 | 说明 | 返回 |
|------|------|------|
| `Check(ident, code)` | 验证验证码是否正确 | `error` |
| `CheckInTime(ident, code, validtime_in_second)` | 验证验证码，并限制有效时间 | `error` |
| `CheckWithCode(ident, code)` | 验证验证码，返回业务状态码 | `(int64, error)` |
| `Math(ident)` | 生成数学算式验证码图片 | `(image.Image, error)` |
| `Number(ident)` | 生成数字验证码图片 | `(image.Image, error)` |
| `Chinese(ident)` | 生成中文验证码图片 | `(image.Image, error)` |
| `Text(ident)` | 生成文本验证码图片（字母+数字） | `(image.Image, error)` |

> **ident 参数说明**：验证码标识，建议用 session id 或用户唯一标识。**每次调用必须使用新的 ident，防止验证码被复用**。

### Check

验证验证码是否正确，返回 `nil` 表示验证通过。

```go
func (self *Captcha) Check(ident, code any) error
```

**使用示例：**

```go
captcha := &AossGoSdk.Captcha{Token: "your-project-token"}

err := captcha.Check("user_session_id", "1234")
if err != nil {
    fmt.Println("验证码错误:", err)
    return
}
// 验证通过
```

### CheckInTime

```go
func (self *Captcha) CheckInTime(ident, code any, validtime_in_second int) error
```

```go
err := captcha.CheckInTime("user_session_id", "1234", 300)
if err != nil {
    fmt.Println("验证码无效或已过期:", err)
    return
}
```

### CheckWithCode

```go
func (self *Captcha) CheckWithCode(ident, code any) (int64, error)
```

验证码不通过时 code 为 `-103`。

```go
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

### 生成图片验证码

以下方法均返回 `image.Image` 类型，可直接写入 HTTP 响应或保存为文件。

**完整示例（net/http）：**

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
        http.Error(w, err.Error(), http.StatusInternalServerError)
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

func main() {
    http.HandleFunc("/captcha", GetCaptcha)
    http.HandleFunc("/captcha/verify", VerifyCaptcha)
    http.ListenAndServe(":8080", nil)
}
```

---

## 二、动态 GIF 验证码

GIF 验证码返回 `[]byte` 类型的 GIF 二进制数据，适合直接写入 HTTP 响应。

| 方法 | 说明 |
|------|------|
| `GifText(ident)` | 字母+数字，1 秒一帧 |
| `GifFast(ident)` | 字母+数字，0.5 秒一帧 |
| `GifNumber(ident)` | 纯数字，1 秒一帧 |
| `GifNumberFast(ident)` | 纯数字，0.5 秒一帧 |
| `GifLetters(ident)` | 纯字母，1 秒一帧 |
| `GifLettersFast(ident)` | 纯字母，0.5 秒一帧 |

**使用示例：**

```go
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

---

## 三、滑动拼图验证码（推荐）

滑动拼图验证码需要前后端配合：**后端负责生成和校验，前端负责渲染图片和滑块交互。**

### 数据结构

```go
type SlideCaptchaData struct {
    Bg        string `json:"bg"`        // 背景图（带缺口），base64 data URL
    Block     string `json:"block"`     // 滑块图（拼图块），base64 data URL
    Y         int    `json:"y"`         // 滑块纵坐标
    BgWidth   int    `json:"bg_width"`  // 背景图宽度
    BgHeight  int    `json:"bg_height"` // 背景图高度
    BlockSize int    `json:"block_size"`// 滑块边长
}
```

### 生成验证码

```go
func (self *Captcha) Slide(ident any) (data SlideCaptchaData, err error)
```

### 验证验证码

```go
func (self *Captcha) SlideCheck(ident any, x int) error
```

```go
func (self *Captcha) SlideCheckWithCode(ident any, x int) (int64, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| ident | any | 与生成时一致的标识 |
| x | int | 滑块最终水平位置（像素） |

---

### 前后端完整示例：Go Gin + 原生 HTML/JS

**后端 `main.go`：**

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

var captcha = &AossGoSdk.Captcha{Token: "your-project-token"}

func main() {
    r := gin.Default()

    // 生成滑动拼图
    r.POST("/captcha/slide/create", func(c *gin.Context) {
        ident := "slide_" + randomString(16)
        data, err := captcha.Slide(ident)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 500,
                "echo": err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code":  0,
            "ident": ident,
            "data":  data,
        })
    })

    // 校验滑动拼图
    r.POST("/captcha/slide/check", func(c *gin.Context) {
        ident := c.PostForm("ident")
        x := c.GetInt("x")
        err := captcha.SlideCheck(ident, x)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": -103,
                "echo": err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code": 0,
            "echo": "验证成功",
        })
    })

    r.StaticFile("/slide", "./slide.html")
    r.Run(":8080")
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
```

**前端页面 `slide.html`：**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>滑动拼图验证码</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
        .captcha-container { position: relative; width: 300px; height: 190px; border: 1px solid #ccc; background: white; box-shadow: 0 0 10px rgba(0,0,0,0.1); overflow: hidden; }
        .captcha-bg { width: 300px; height: 150px; display: block; position: absolute; top: 0; left: 0; }
        .captcha-block { position: absolute; left: 0; top: 0; cursor: grab; z-index: 10; transition: opacity 0.3s ease; }
        .captcha-block:active { cursor: grabbing; }
        .captcha-slider { position: absolute; bottom: 0; left: 0; width: 100%; height: 40px; background: #f5f5f5; border-top: 1px solid #eee; }
        .slider-handle { width: 40px; height: 100%; background: #409eff; color: white; text-align: center; line-height: 40px; cursor: grab; position: absolute; font-size: 16px; border-radius: 2px; transition: background-color 0.2s ease; user-select: none; }
        .slider-handle:active { cursor: grabbing; background: #66b1ff; }
        .slider-handle.success { background: #67c23a; }
        .slider-handle.error { background: #f56c6c; }
        .slide-status { margin-top: 12px; text-align: center; color: #666; min-height: 24px; }
        .slide-status.success { color: #67c23a; font-weight: bold; }
        .slide-status.error { color: #f56c6c; font-weight: bold; }
    </style>
</head>
<body>
    <div>
        <div class="captcha-container">
            <img class="captcha-bg" src="" alt="验证码背景">
            <img class="captcha-block" src="" alt="验证码块" style="display:none;">
            <div class="captcha-slider"><div class="slider-handle">👉</div></div>
        </div>
        <div class="slide-status">请拖动下方滑块将拼图拼合</div>
    </div>

    <script>
    let captchaData = null;
    let captchaIdent = '';
    let startX = 0;
    let startLeft = 0;

    const bgImg = document.querySelector('.captcha-bg');
    const blockImg = document.querySelector('.captcha-block');
    const sliderHandle = document.querySelector('.slider-handle');
    const statusEl = document.querySelector('.slide-status');

    function generate() {
        fetch('/captcha/slide/create', { method: 'POST' })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    captchaData = data.data;
                    captchaIdent = data.ident;
                    captchaData.pad_top = captchaData.pad_top || 0;
                    captchaData.pad_left = captchaData.pad_left || 0;

                    bgImg.src = captchaData.bg;
                    blockImg.style.width = captchaData.block_size + 'px';
                    blockImg.style.height = captchaData.block_size + 'px';
                    blockImg.style.top = (captchaData.y - captchaData.pad_top) + 'px';
                    blockImg.style.left = (-captchaData.pad_left) + 'px';
                    blockImg.onload = function() { blockImg.style.display = 'block'; };
                    blockImg.src = captchaData.block;

                    sliderHandle.style.left = '0px';
                    sliderHandle.classList.remove('success', 'error');
                    statusEl.textContent = '请拖动下方滑块将拼图拼合';
                    statusEl.className = 'slide-status';
                } else {
                    statusEl.textContent = data.echo || '生成验证码失败';
                    statusEl.className = 'slide-status error';
                }
            });
    }

    function startDrag(e) {
        if (!captchaData) return;
        startX = e.clientX || (e.touches && e.touches[0].clientX);
        startLeft = parseInt(sliderHandle.style.left) || 0;
    }

    function drag(e) {
        if (!captchaData) return;
        const clientX = e.clientX || (e.touches && e.touches[0].clientX);
        const deltaX = clientX - startX;
        const newLeft = Math.max(0, Math.min(startLeft + deltaX, captchaData.bg_width - captchaData.block_size));
        sliderHandle.style.left = newLeft + 'px';
        blockImg.style.left = (newLeft - captchaData.pad_left) + 'px';
    }

    function endDrag(e) {
        if (!captchaData) return;
        const finalX = parseInt(sliderHandle.style.left) || 0;

        const formData = new FormData();
        formData.append('ident', captchaIdent);
        formData.append('x', finalX);

        fetch('/captcha/slide/check', { method: 'POST', body: formData })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    sliderHandle.classList.add('success');
                    statusEl.textContent = '✅ 验证成功！';
                    statusEl.className = 'slide-status success';
                } else {
                    sliderHandle.classList.add('error');
                    statusEl.textContent = '❌ ' + (data.echo || '验证失败') + '，正在刷新...';
                    statusEl.className = 'slide-status error';
                    setTimeout(generate, 1500);
                }
            });
    }

    sliderHandle.addEventListener('mousedown', startDrag);
    document.addEventListener('mousemove', drag);
    document.addEventListener('mouseup', endDrag);
    sliderHandle.addEventListener('touchstart', startDrag, { passive: true });
    document.addEventListener('touchmove', drag, { passive: true });
    document.addEventListener('touchend', endDrag);

    generate();
    </script>
</body>
</html>
```

---

## 四、点击验证码（安全级别最高）

点击验证码要求用户根据提示**依次点击图片中的特定图标**，需前后端配合。

### 数据结构

```go
type ClickCaptchaData struct {
    Bg           string `json:"bg"`            // 背景图（带有多个可点击图标），base64 data URL
    Tip          string `json:"tip"`           // 给用户看的提示文字，如 "请依次点击：A、B、C"
    TargetsCount int    `json:"targets_count"` // 需要点击的图标数量
    BgWidth      int    `json:"bg_width"`      // 背景图宽度
    BgHeight     int    `json:"bg_height"`     // 背景图高度
}
```

### 生成验证码

```go
func (self *Captcha) Click(ident any) (data ClickCaptchaData, err error)
```

### 验证验证码

```go
func (self *Captcha) ClickCheck(ident any, clicks string) error
```

| 参数 | 类型 | 说明 |
|------|------|------|
| ident | any | 与生成时一致的标识 |
| clicks | string | 点击坐标 JSON 字符串，格式 `[{"x":100,"y":80}, {"x":200,"y":120}]` |

---

### 前后端完整示例：Go Gin + 原生 HTML/JS

**后端 `main.go`：**

```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/gin-gonic/gin"
    AossGoSdk "github.com/tobycroft/AossGoSdk"
)

type ClickPoint struct {
    X int `json:"x"`
    Y int `json:"y"`
}

var captcha = &AossGoSdk.Captcha{Token: "your-project-token"}

func main() {
    r := gin.Default()

    // 生成点击验证码
    r.POST("/captcha/click/create", func(c *gin.Context) {
        ident := "click_" + randomString(16)
        data, err := captcha.Click(ident)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "code": 500,
                "echo": err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code":  0,
            "ident": ident,
            "data":  data,
        })
    })

    // 校验点击验证码
    r.POST("/captcha/click/check", func(c *gin.Context) {
        ident := c.PostForm("ident")
        clicksJson := c.PostForm("clicks")
        err := captcha.ClickCheck(ident, clicksJson)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code": -103,
                "echo": err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "code": 0,
            "echo": "验证成功",
        })
    })

    r.StaticFile("/click", "./click.html")
    r.Run(":8080")
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
```

**前端页面 `click.html`：**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>点击验证码</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f5f5f5; }
        .click-captcha-wrapper { text-align: center; }
        .click-captcha-container { position: relative; width: 300px; height: 200px; border: 1px solid #ccc; background: white; box-shadow: 0 0 10px rgba(0,0,0,0.1); cursor: default; display: inline-block; user-select: none; -webkit-user-select: none; }
        .click-captcha-bg { width: 300px; height: 200px; display: block; }
        .click-marker { position: absolute; width: 22px; height: 22px; line-height: 22px; text-align: center; background: rgba(64, 158, 255, 0.8); color: white; border-radius: 50%; font-size: 12px; font-weight: bold; transform: translate(-50%, -50%); pointer-events: none; z-index: 10; animation: marker-pop 0.2s ease-out; }
        @keyframes marker-pop { 0% { transform: translate(-50%, -50%) scale(0.5); opacity: 0; } 100% { transform: translate(-50%, -50%) scale(1); opacity: 1; } }
        .click-tip { margin-top: 12px; font-size: 15px; color: #333; font-weight: bold; min-height: 22px; }
        .click-count { margin-top: 4px; font-size: 13px; color: #999; min-height: 20px; }
        .click-status { margin-top: 8px; font-size: 14px; min-height: 20px; transition: all 0.3s ease; }
        .click-status.success { color: #67c23a; font-weight: bold; }
        .click-status.error { color: #f56c6c; font-weight: bold; }
        .click-reload-btn { margin-top: 10px; padding: 8px 16px; background: #409eff; color: white; border: none; border-radius: 4px; cursor: pointer; transition: background 0.3s ease; }
        .click-reload-btn:hover { background: #66b1ff; }
        .click-reload-btn:disabled { background: #ccc; cursor: not-allowed; }
    </style>
</head>
<body>
    <div class="click-captcha-wrapper">
        <div class="click-captcha-container">
            <img class="click-captcha-bg" src="" alt="验证码背景">
        </div>
        <div class="click-tip"></div>
        <div class="click-count"></div>
        <div class="click-status"></div>
        <button class="click-reload-btn">刷新验证码</button>
    </div>

    <script>
    let captchaData = null;
    let captchaIdent = '';
    let userClicks = [];
    let targetCount = 0;

    const bgImg = document.querySelector('.click-captcha-bg');
    const tipEl = document.querySelector('.click-tip');
    const clickCountEl = document.querySelector('.click-count');
    const statusEl = document.querySelector('.click-status');
    const reloadBtn = document.querySelector('.click-reload-btn');
    const container = document.querySelector('.click-captcha-container');

    function generate() {
        userClicks = [];
        const markers = container.querySelectorAll('.click-marker');
        markers.forEach(m => m.remove());

        fetch('/captcha/click/create', { method: 'POST' })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    captchaData = data.data;
                    captchaIdent = data.ident;
                    targetCount = captchaData.targets_count;
                    bgImg.src = captchaData.bg;
                    tipEl.textContent = captchaData.tip;
                    clickCountEl.textContent = '还需点击 ' + targetCount + ' 个';
                    statusEl.textContent = '';
                    statusEl.className = 'click-status';
                } else {
                    statusEl.textContent = data.echo || '生成验证码失败';
                    statusEl.className = 'click-status error';
                }
            });
    }

    function drawMarker(x, y, num) {
        const marker = document.createElement('div');
        marker.className = 'click-marker';
        marker.style.left = x + 'px';
        marker.style.top = y + 'px';
        marker.textContent = num;
        container.appendChild(marker);
    }

    function verify() {
        const formData = new FormData();
        formData.append('ident', captchaIdent);
        formData.append('clicks', JSON.stringify(userClicks));

        fetch('/captcha/click/check', { method: 'POST', body: formData })
            .then(r => r.json())
            .then(data => {
                if (data.code === 0) {
                    statusEl.textContent = '✅ 验证成功！';
                    statusEl.className = 'click-status success';
                } else {
                    statusEl.textContent = '❌ ' + (data.echo || '验证失败') + '，正在刷新...';
                    statusEl.className = 'click-status error';
                    setTimeout(generate, 1500);
                }
            });
    }

    container.addEventListener('click', function(e) {
        if (!captchaData) return;
        const rect = container.getBoundingClientRect();
        const x = Math.round(e.clientX - rect.left);
        const y = Math.round(e.clientY - rect.top);
        if (y >= captchaData.bg_height) return;

        userClicks.push({ x: x, y: y });
        drawMarker(x, y, userClicks.length);

        const remaining = targetCount - userClicks.length;
        if (remaining > 0) {
            clickCountEl.textContent = '还需点击 ' + remaining + ' 个';
        } else {
            clickCountEl.textContent = '正在验证...';
            verify();
        }
    });

    reloadBtn.addEventListener('click', function() {
        if (reloadBtn.disabled) return;
        reloadBtn.disabled = true;
        generate();
        setTimeout(function() { reloadBtn.disabled = false; }, 500);
    });

    generate();
    </script>
</body>
</html>
```

---

## 五、通用接口参数

### 文本/GIF 验证码

| 步骤 | 接口 | 参数 | 返回 |
|------|------|------|------|
| 生成 | `POST /v1/captcha/text/{type}` | `token`, `ident` | PNG 图片二进制 / GIF 二进制 |
| 校验 | `POST /v1/captcha/auth/check` | `token`, `ident`, `code` | `{ code, echo }` |
| 校验（带有效期） | `POST /v1/captcha/auth/check_in_time` | `token`, `ident`, `code`, `second` | `{ code, echo }` |

### 滑动拼图验证码

| 步骤 | 接口 | 参数 | 返回 |
|------|------|------|------|
| 生成 | `POST /v1/captcha/slide/create` | `token`, `ident` | `{ code, data: { bg, block, y, bg_width, bg_height, block_size, pad_top, pad_left } }` |
| 校验 | `POST /v1/captcha/slide/check` | `token`, `ident`, `x` | `{ code, echo }` |

### 点击验证码

| 步骤 | 接口 | 参数 | 返回 |
|------|------|------|------|
| 生成 | `POST /v1/captcha/click/create` | `token`, `ident` | `{ code, data: { bg, tip, targets_count, bg_width, bg_height } }` |
| 校验 | `POST /v1/captcha/click/check` | `token`, `ident`, `clicks`（JSON 字符串） | `{ code, echo }` |

---

## 注意事项

1. **Ident 唯一**：每次调用 `Slide` / `Click` / 生成图片时，务必使用新的 ident，防止验证码被复用攻击
2. **Token 安全**：生产环境建议 token 由后端管理，不要直接暴露在前端 JS 中；可通过后端代理转发请求
3. **响应式适配**：上述前端示例中 `.captcha-container` 固定为 300px。如需移动端适配，可改为 `width: 100%` 并动态调整尺寸
4. **触摸支持**：示例已内置 `touchstart/touchmove/touchend`，可直接在移动端使用
5. **HTTP 缓存**：验证码接口返回默认禁止缓存，前端 fetch 时请不要添加缓存头