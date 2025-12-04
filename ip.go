package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
	Net "github.com/tobycroft/TuuzNet"
)

type IP struct {
	Code  string
	Token string
}

func (self *IP) IpRange(country any, province []any, ip any) (b bool, err error) {
	province_json, err := sonic.MarshalString(province)
	if err != nil {
		return false, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/ip/range/check", nil, param, nil, nil)
	//fmt.Println(ret, err)
	var rs ret_std
	err = rets.RetJson(&rs)
	if err != nil {
		return
	} else {
		if rs.Code == 0 {
			return true, nil
		} else {
			return false, errors.New(rs.Echo)
		}
	}
}

// IpRangeAuth this will check the ip is in the ip range of the country and province,if not it will needs client to complete the captcha check, this is recommended to be used if you are running a sms services
func (self *IP) IpRangeAuth(country any, province []any, ip any) (bool, error) {
	province_json, err := sonic.MarshalString(province)
	if err != nil {
		return false, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/ip/range/auth", nil, param, nil, nil)
	//fmt.Println(ret, err)
	var rs ret_std
	err = rets.RetJson(&rs)
	if err != nil {
		return false, err
	}
	if rs.Code == 0 {
		return true, nil
	} else {
		return false, errors.New(rs.Echo)
	}
}

// IpRangeCaptcha This function is pretty similar to IpRangeAuth, but it will return the response code like CheckWithCode, when the code is -103 means it needs to complete the captcha check
func (self *IP) IpRangeCaptcha(country any, province []any, ip any) (int64, error) {
	province_json, err := sonic.MarshalString(province)
	if err != nil {
		return 400, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+"/v1/ip/range/auth", nil, param, nil, nil)
	//fmt.Println(ret, err)
	var rs ret_std
	err = rets.RetJson(&rs)
	if err != nil {
		return 500, err
	} else {
		if rs.Code == 0 {
			return rs.Code, nil
		} else {
			return rs.Code, errors.New(rs.Echo)
		}
	}
}
