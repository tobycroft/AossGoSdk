package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

const Canvas_Type_text = "text"
const Canvas_Type_image = "image"
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
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type Canvas_Type_Image struct {
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type Canvas_type struct {
	Canvas_Type_Text
	Canvas_Type_Image
}

// Canvas_url:获取微信小程序二维码（302方法，推荐占用少）
func Canvas_url(project, width int64, height int64, background string, canvas_type []Canvas_type) (string, error) {
	data, err := jsoniter.MarshalToString(canvas_type)
	if err != nil {
		return "", err
	}
	post := map[string]any{
		"width":      width,
		"height":     height,
		"background": background,
		"data":       data,
	}
	ret, err := Net.Post(baseUrl+canvas_file, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return "", err
	}
	var wwuf wechatStringData
	err = jsoniter.UnmarshalFromString(ret, &wwuf)
	if err != nil {
		return "", err
	}
	if wwuf.Code == 0 {
		return wwuf.Data, nil
	} else {
		return wwuf.Data, errors.New(wwuf.Echo)
	}
}
