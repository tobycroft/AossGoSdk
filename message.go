package AossGoSdk

import Net "github.com/tobycroft/TuuzNet"

type Wechat_message struct {
	Token    string
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

// set_token:如果不设定token
func (self *Wechat_message) set_token(token string) {

}

func (self *Wechat_message) Send() {
	post := map[string]any{
		"openid":  self.openid,
		"content": self.content,
	}
	ret, err := Net.Post(baseUrls+message_custom_send, map[string]interface{}{
		"token": self.Token,
	}, post, nil, nil)
}
