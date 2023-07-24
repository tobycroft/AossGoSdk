package AossGoSdk

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
)

type ASMS struct {
	Name  string
	Token string
}

func (self *ASMS) Sms_send(phone any, quhao, text, ip any) error {
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
	ret, err := Net.Post(baseUrl+"/v1/sms/single/push", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return err
	} else {
		var rs ret_std
		errs := json.Unmarshal([]byte(ret), &rs)
		if errs != nil {
			return errors.New(ret)
		} else {
			//fmt.Println(rs)
			if rs.Code == 0 {
				return nil
			} else {
				return errors.New(rs.Echo)
			}
		}
	}
}
