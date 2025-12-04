package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
)

const Canvas_Posistion_TopLeft = "lt"
const Canvas_Posistion_TopCenter = "mt"
const Canvas_Posistion_TopRight = "rt"
const Canvas_Posistion_CenterLeft = "lm"
const Canvas_Posistion_CenterCenter = "mm"
const Canvas_Posistion_CenterRight = "rm"
const Canvas_Posistion_BottomLeft = "lb"
const Canvas_Posistion_BottomCenter = "mb"
const Canvas_Posistion_BottomRight = "rb"

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

type Canvas struct {
	layer []interface{}
}

func (self *Canvas) AddText(Text string, Canvas_position string, X int64, Y int64) *Canvas {
	self.layer = append(self.layer, Canvas_Type_Text{
		Type:     "text",
		Text:     Text,
		Position: Canvas_position,
		X:        X,
		Y:        Y,
	})
	return self
}

func (self *Canvas) AddImage(Url string, X int64, Y int64) *Canvas {
	self.layer = append(self.layer, Canvas_Type_Image{
		Type: "image",
		URL:  Url,
		X:    X,
		Y:    Y,
	})
	return self
}

// Canvas_url:获取微信小程序二维码（302方法，推荐占用少）
func (self *Canvas) Get_Url(project interface{}, width int64, height int64, background_color string) (str string, err error) {
	data, err := sonic.MarshalString(self.layer)
	if err != nil {
		return "", err
	}
	post := map[string]string{
		"width":      Calc.Any2String(width),
		"height":     Calc.Any2String(height),
		"background": background_color,
		"data":       data,
	}
	var wwuf wechatStringData
	err2 := new(Net.Post).PostFormData(baseUrl+canvas_file, map[string]interface{}{"token": project}, post, nil, nil).RetJson(&wwuf)
	if err2 != nil {
		return "", err2
	}
	if wwuf.Code == 0 {
		return wwuf.Data, nil
	} else {
		return wwuf.Data, errors.New(wwuf.Echo)
	}
}
