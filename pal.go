package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

type PalWorld struct {
	Name  string
	Token string
}
type OnlineUser struct {
	Name      string `json:"name"`
	Playeruid string `json:"playeruid"`
	Steamid   string `json:"steamid"`
}
type onlineUser struct {
	Code int          `json:"code"`
	Data []OnlineUser `json:"data"`
	Echo string       `json:"echo"`
}

// Lcic_CreateUser Name:用户在直播间的显示名称 | OriginId:用户在你系统中的标识符 | Avatar:用户头像url地址
func (self *Lcic) OnlineUser() ([]OnlineUser, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/rcon/palworld/players", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return nil, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return nil, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat onlineUser
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return nil, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return nil, errors.New(resp.Echo)
		}
	}
}
