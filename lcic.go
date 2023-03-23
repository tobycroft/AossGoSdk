package AossGoSdk

import (
	"encoding/json"
	"errors"
	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

type Lcic struct {
	Name  string
	Token string
}

func (self *Lcic) Lcic_CreateUser(Name, OriginId, Avatar interface{}) error {
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

func Lcic_RoomCreate(TeacherId, StartTime, EndTime, Name interface{}) {

}

func Lcic_RoomModify(RoomId, TeacherId, StartTime, EndTime, Name interface{}) {

}

func Lcic_RoomDelete(RoomId interface{}) {

}
