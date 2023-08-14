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
