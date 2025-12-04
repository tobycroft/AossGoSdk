package AossGoSdk

import (
	"errors"
	"time"

	Net "github.com/tobycroft/TuuzNet"
)

// Wechat_ticket_signature:微信ticket签名功能
func Wechat_ticket_signature(project, noncestr string, timestamp time.Time, url string) (str string, err error) {
	post := map[string]any{
		"noncestr":  noncestr,
		"timestamp": timestamp.Unix(),
		"url":       url,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+ticket_signature, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return "", err
	}

	if resp.Code == 0 {
		var wwuf wechatStringData
		err = rets.RetJson(&wwuf)
		if err != nil {
			return "", err
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}
