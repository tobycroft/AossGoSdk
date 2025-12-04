package AossGoSdk

import (
	"errors"
	"time"

	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
)

type ASMS struct {
	Name  string
	Token string
}

func (self *ASMS) Sms_send(phone any, quhao, text, ip any) (err error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"phone": phone,
		"quhao": quhao,
		"text":  text,
		"ip":    ip,
		"ts":    ts,
		"name":  self.Name,
		"sign":  Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/sms/single/push", nil, param, nil, nil)
	var rs ret_std
	err = rets.RetJson(&rs)
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}
