package AossGoSdk

import (
	"encoding/json"
	"errors"
	Net "github.com/tobycroft/TuuzNet"
)

type Captcha struct {
	Token any
}

func (self *Captcha) Check(ident, code any) error {
	param := map[string]any{
		"ident": ident,
		"code":  code,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/check", nil, param, nil, nil)
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

func (self *Captcha) Math(ident any) (string, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/math", nil, param, nil, nil)
	//fmt.Println(ret, err)
	return ret, err
}
func (self *Captcha) Number(ident any) (string, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/number", nil, param, nil, nil)
	//fmt.Println(ret, err)
	return ret, err
}
func (self *Captcha) Chinese(ident any) (string, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/chinese", nil, param, nil, nil)
	//fmt.Println(ret, err)
	return ret, err
}
func (self *Captcha) Text(ident any) (string, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/text", nil, param, nil, nil)
	//fmt.Println(ret, err)
	return ret, err
}
