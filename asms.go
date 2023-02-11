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

func (self *ASMS) Sms_send(phone any, quhao, text any) error {
	ts := time.Now().Unix()
	param := map[string]any{
		"phone": phone,
		"quhao": quhao,
		"text":  text,
		"ts":    ts,
		"name":  self.Name,
		"sign":  Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/sms/single/push", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return err
	} else {
		var rs ret_struct
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

type ret_struct struct {
	Code int64
	Echo string
}
