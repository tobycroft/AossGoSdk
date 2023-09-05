package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

type IP struct {
	Token string
	Code  string
}

func (self *IP) IpRange(country any, province []any, ip any) (bool, error) {
	province_json, err := jsoniter.MarshalToString(province)
	if err != nil {
		return false, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/ip/range/check", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return false, err
	} else {
		var rs ret_std
		errs := jsoniter.UnmarshalFromString(ret, &rs)
		if errs != nil {
			return false, errors.New(ret)
		} else {
			if rs.Code == 0 {
				return true, nil
			} else {
				return false, errors.New(rs.Echo)
			}
		}
	}
}

// IpRangeAuth this will check the ip is in the ip range of the country and province,if not it will needs client to complete the captcha check, this is recommended to be used if you are running a sms services
func (self *IP) IpRangeAuth(country any, province []any, ip any) (bool, error) {
	province_json, err := jsoniter.MarshalToString(province)
	if err != nil {
		return false, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/ip/range/auth", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return false, err
	} else {
		var rs ret_std
		errs := jsoniter.UnmarshalFromString(ret, &rs)
		if errs != nil {
			return false, errors.New(ret)
		} else {
			if rs.Code == 0 {
				return true, nil
			} else {
				return false, errors.New(rs.Echo)
			}
		}
	}
}

// IpRangeCaptcha This function is pretty similar to IpRangeAuth, but it will return the response code like CheckWithCode, when the code is -103 means it needs to complete the captcha check
func (self *IP) IpRangeCaptcha(country any, province []any, ip any) (int64, error) {
	province_json, err := jsoniter.MarshalToString(province)
	if err != nil {
		return 400, err
	}
	param := map[string]any{
		"country":  country,
		"province": province_json,
		"ip":       ip,
		"token":    self.Token,
	}
	ret, err := Net.Post(baseUrl+"/v1/ip/range/auth", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return 200, err
	} else {
		var rs ret_std
		errs := jsoniter.UnmarshalFromString(ret, &rs)
		if errs != nil {
			return 500, errors.New(ret)
		} else {
			if rs.Code == 0 {
				return rs.Code, nil
			} else {
				return rs.Code, errors.New(rs.Echo)
			}
		}
	}
}
