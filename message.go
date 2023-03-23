package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

type Wechat_message struct {
	token    string
	msg_type string
	openid   string
	content  string
}

func (self *Wechat_message) Set_message_text(content string) *Wechat_message {
	self.msg_type = "text"
	self.content = content
	return self
}

func (self *Wechat_message) Set_openid(openid string) *Wechat_message {
	self.openid = openid
	return self
}

// set_token:如果不设定token
func (self *Wechat_message) Set_token_not_set_for_auto(token string) *Wechat_message {
	self.token = token
	return self
}

func (self *Wechat_message) Send() error {
	post := map[string]any{
		"openid":  self.openid,
		"content": self.content,
	}
	get := map[string]interface{}{
		"token": self.token,
	}
	url := ""
	switch self.msg_type {
	case "text":
		url = baseUrls + message_send_text
		break

	default:
		url = baseUrls + message_send_text
		break
	}
	ret, err := Net.Post(url, get, post, nil, nil)
	if err != nil {
		return err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return err
		}
		if resp.Code == 0 {
			return nil
		} else {
			return errors.New(resp.Echo)
		}
	}
}

type Wechat_message_ret_struct struct {
	Project      string
	ToUserName   string
	FromUserName string
	CreateTime   any
	MsgType      string
	Event        string
	EventKey     any
	Ticket       string
	Content      string
	MsgId        any
	Idx          string
}
