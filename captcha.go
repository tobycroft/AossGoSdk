package AossGoSdk

import (
	"bytes"
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
	"image"
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
	ret, err := Net.Post(baseUrl+"/v1/captcha/auth/check", nil, param, nil, nil)
	if err != nil {
		return err
	} else {
		var rs ret_std
		errs := jsoniter.UnmarshalFromString(ret, &rs)
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

// CheckWithCode This function will return the httpcode if the captcha is not valid, the code will be returned as -103, code 200 for network error, code 500 for json decode error
func (self *Captcha) CheckWithCode(ident, code any) (int64, error) {
	param := map[string]any{
		"ident": ident,
		"code":  code,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/auth/check", nil, param, nil, nil)
	if err != nil {
		return 200, err
	} else {
		var rs ret_std
		errs := jsoniter.UnmarshalFromString(ret, &rs)
		if errs != nil {
			return 500, errors.New(ret)
		} else {
			//fmt.Println(rs)
			if rs.Code == 0 {
				return rs.Code, nil
			} else {
				return rs.Code, errors.New(rs.Echo)
			}
		}
	}
}

func (self *Captcha) Math(ident any) (image.Image, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/math", nil, param, nil, nil)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader([]byte(ret)))
	return img, err
}
func (self *Captcha) Number(ident any) (image.Image, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/number", nil, param, nil, nil)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader([]byte(ret)))
	return img, err
}
func (self *Captcha) Chinese(ident any) (image.Image, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/chinese", nil, param, nil, nil)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader([]byte(ret)))
	return img, err
}
func (self *Captcha) Text(ident any) (image.Image, error) {
	param := map[string]any{
		"ident": ident,
		"token": self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/captcha/text/text", nil, param, nil, nil)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader([]byte(ret)))
	return img, err
}
