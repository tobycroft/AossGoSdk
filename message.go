package AossGoSdk

import Net "github.com/tobycroft/TuuzNet"

type Wechat_message struct {
	msg_type string
	openid   string
	content  string
}

func (self *Wechat_message) set_message_type_text() *Wechat_message {
	self.msg_type = "text"
	return self
}

func (self *Wechat_message) set_openid(openid string) *Wechat_message {
	self.openid = openid
	return self
}

func (self *Wechat_message) Send() {
	post := map[string]any{
		"openid":  self.openid,
		"content": self.content,
	}
	ret, err := Net.Post(baseUrls+message_custom_send, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
}
