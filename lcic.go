package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

type Lcic struct {
	Name  string
	Token string
}

type lcicStructCreateUser struct {
	Code int                  `json:"code"`
	Data LcicStructCreateUser `json:"data"`
	Echo string               `json:"echo"`
}

type LcicStructCreateUser struct {
	UserId string `json:"UserId"`
	Token  string `json:"Token"`
}

func (self *Lcic) Lcic_CreateUser(Name, OriginId, Avatar interface{}) (LcicStructCreateUser, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"Name":     Name,
		"OriginId": OriginId,
		"Avatar":   Avatar,
		"ts":       ts,
		"name":     self.Name,
		"sign":     Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/sms/single/push", nil, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LcicStructCreateUser{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return LcicStructCreateUser{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat lcicStructCreateUser
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LcicStructCreateUser{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LcicStructCreateUser{}, errors.New(resp.Echo)
		}
	}
}

func Lcic_RoomCreate(TeacherId, StartTime, EndTime, Name interface{}) {

}

func Lcic_RoomModify(RoomId, TeacherId, StartTime, EndTime, Name interface{}) {

}

func Lcic_RoomDelete(RoomId interface{}) {

}
