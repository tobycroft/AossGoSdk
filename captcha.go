package AossGoSdk

import (
	"bytes"
	"errors"
	"image"

	Net "github.com/tobycroft/TuuzNet"
)

type Captcha struct {
	Token any
}

// Check if the error returned not-nil that represents the captcha has errors
func (self *Captcha) Check(ident, code any) error {
	param := map[string]any{
		"ident": ident,
		"code":  code,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/auth/check", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return err
	}
	//fmt.Println(rs)
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}

// Check if the error returned not-nil that represents the captcha has errors
func (self *Captcha) CheckInTime(ident, code any, validtime_in_second int) error {
	param := map[string]any{
		"ident":  ident,
		"code":   code,
		"second": validtime_in_second,
		"token":  self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/auth/check_in_time", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return err
	}
	//fmt.Println(rs)
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}

// CheckWithCode This function will return the httpcode if the captcha is not valid, the code will be returned as -103, code 200 for network error, code 500 for json decode error
func (self *Captcha) CheckWithCode(ident, code any) (int64, error) {
	param := map[string]any{
		"ident": ident,
		"code":  code,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/auth/check", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return 200, err
	}
	if rs.Code == 0 {
		return rs.Code, nil
	} else {
		return rs.Code, errors.New(rs.Echo)
	}
}

func (self *Captcha) Math(ident any) (img image.Image, err error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/text/math", nil, param, nil, nil)
	var ret []byte
	ret, err = rets.RetBytes()
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(ret))
	if err != nil {
		var rs ret_std
		err = rets.RetJson(&rs)
		if err != nil {
			return
		}
		//fmt.Println(rs)
		if rs.Code == 0 {
			return
		} else {
			err = errors.New(rs.Echo)
			return
		}
	}
	return img, err
}
func (self *Captcha) Number(ident any) (img image.Image, err error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/text/number", nil, param, nil, nil)
	var ret []byte
	ret, err = rets.RetBytes()
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(ret))
	if err != nil {
		var rs ret_std
		err = rets.RetJson(&rs)
		if err != nil {
			return
		}
		//fmt.Println(rs)
		if rs.Code == 0 {
			return
		} else {
			err = errors.New(rs.Echo)
			return
		}
	}
	return img, err
}

func (self *Captcha) Chinese(ident any) (img image.Image, err error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/text/chinese", nil, param, nil, nil)
	var ret []byte
	ret, err = rets.RetBytes()
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(ret))
	if err != nil {
		var rs ret_std
		err = rets.RetJson(&rs)
		if err != nil {
			return
		}
		//fmt.Println(rs)
		if rs.Code == 0 {
			return
		} else {
			err = errors.New(rs.Echo)
			return
		}
	}
	return img, err
}

func (self *Captcha) Text(ident any) (img image.Image, err error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/captcha/text/text", nil, param, nil, nil)
	var ret []byte
	ret, err = rets.RetBytes()
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(ret))
	if err != nil {
		var rs ret_std
		err = rets.RetJson(&rs)
		if err != nil {
			return
		}
		//fmt.Println(rs)
		if rs.Code == 0 {
			return
		} else {
			err = errors.New(rs.Echo)
			return
		}
	}
	return img, err
}
