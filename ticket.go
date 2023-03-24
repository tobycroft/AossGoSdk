package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

// Wechat_ticket_signature:微信ticket签名功能
func Wechat_ticket_signature(project, noncestr string, timestamp time.Time, url string) (string, error) {
	post := map[string]any{
		"noncestr":  noncestr,
		"timestamp": timestamp.Unix(),
		"url":       url,
	}
	ret, err := Net.Post(baseUrl+ticket_signature, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return "", err
	}
	var resp ret_std
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return "", err
	}

	if resp.Code == 0 {
		var wwuf wechatStringData
		err = jsoniter.UnmarshalFromString(ret, &wwuf)
		if err != nil {
			return "", err
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}
