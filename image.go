package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

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

// Canvas_url:获取微信小程序二维码（302方法，推荐占用少）
func Canvas_url(project, width int64, height int64, background string, data string) (string, error) {
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
