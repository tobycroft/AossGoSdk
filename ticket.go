package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

func Wechat_ticket_signature(project, noncestr string, timestamp time.Time, url string) (string, error) {
	post := map[string]any{
		"noncestr":  noncestr,
		"timestamp": timestamp,
		"url":       url,
	}
	ret, err := Net.Post(baseUrl+ticket_signature, map[string]interface{}{
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
