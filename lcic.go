package AossGoSdk

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/tobycroft/Calc"
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

// Lcic_CreateUser Name:用户在直播间的显示名称 | OriginId:用户在你系统中的标识符 | Avatar:用户头像url地址
func (self *Lcic) Lcic_CreateUser(Name, OriginId, Avatar string) (user LcicStructCreateUser, err error) {
	ts := Calc.Any2String(time.Now().Unix())
	param := map[string]string{
		"Name":     Name,
		"OriginId": OriginId,
		"Avatar":   Avatar,
		"ts":       ts,
		"name":     self.Name,
		"sign":     Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormData(baseUrls+"/v1/lcic/user/auto", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat lcicStructCreateUser
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		user = dat.Data
	} else {
		err = errors.New(resp.Echo)
	}
	return
}

type lcicStructCreateRoom struct {
	Code int                  `json:"code"`
	Data LcicStructCreateRoom `json:"data"`
	Echo string               `json:"echo"`
}
type LcicStructCreateRoom struct {
	RoomId int `json:"RoomId"`
}

// Lcic_RoomCreate TeacherId:老师ID | StartTime:开始时间int | EndTime:结束时间int | Name:房间名称
func (self *Lcic) Lcic_RoomCreate(TeacherId, StartTime, EndTime, Name string) (room LcicStructCreateRoom, err error) {
	ts := Calc.Any2String(time.Now().Unix())
	param := map[string]string{
		"Name":      Name,
		"TeacherId": TeacherId,
		"StartTime": StartTime,
		"EndTime":   EndTime,
		"ts":        ts,
		"name":      self.Name,
		"sign":      Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormData(baseUrls+"/v1/lcic/room/create", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat lcicStructCreateRoom
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return LcicStructCreateRoom{}, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return LcicStructCreateRoom{}, errors.New(resp.Echo)
	}
}

// Lcic_RoomModify RoomId:房间ID | TeacherId:老师ID | StartTime:开始时间int | EndTime:结束时间int | Name:房间名称
func (self *Lcic) Lcic_RoomModify(RoomId, TeacherId, StartTime, EndTime, Name string) (bool, error) {
	ts := Calc.Any2String(time.Now().Unix())
	param := map[string]string{
		"RoomId":    RoomId,
		"Name":      Name,
		"TeacherId": TeacherId,
		"StartTime": StartTime,
		"EndTime":   EndTime,
		"ts":        ts,
		"name":      self.Name,
		"sign":      Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/lcic/room/create", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return false, err
	} else {
		var rs ret_std
		errs := json.Unmarshal([]byte(ret), &rs)
		if errs != nil {
			return false, errs
		} else {
			//fmt.Println(rs)
			if rs.Code == 0 {
				return true, nil
			} else {
				return false, errors.New(rs.Echo)
			}
		}
	}
}

// Lcic_RoomDelete RoomId:房间ID
func (self *Lcic) Lcic_RoomDelete(RoomId interface{}) (bool, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"RoomId": RoomId,
		"ts":     ts,
		"name":   self.Name,
		"sign":   Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/lcic/room/delete", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return false, err
	} else {
		var rs ret_std
		errs := json.Unmarshal([]byte(ret), &rs)
		if errs != nil {
			return false, errs
		} else {
			//fmt.Println(rs)
			if rs.Code == 0 {
				return true, nil
			} else {
				return false, errors.New(rs.Echo)
			}
		}
	}
}

type lcicStructLinkUrl struct {
	Code int               `json:"code"`
	Data LcicStructLinkUrl `json:"data"`
	Echo string            `json:"echo"`
}

type LcicStructLinkUrl struct {
	Web string `json:"web"`
	Pc  string `json:"pc"`
}

// Lcic_LinkUrl OriginId:学生id | TeacherId:老师id
func (self *Lcic) Lcic_LinkUrl(OriginId, TeacherId interface{}) (LcicStructLinkUrl, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"OriginId":  OriginId,
		"TeacherId": TeacherId,
		"ts":        ts,
		"name":      self.Name,
		"sign":      Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/lcic/room/delete", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LcicStructLinkUrl{}, err
	} else {
		var resp ret_std
		err = rets.RetJson(&resp)
		if err != nil {
			return LcicStructLinkUrl{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat lcicStructLinkUrl
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LcicStructLinkUrl{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LcicStructLinkUrl{}, errors.New(resp.Echo)
		}
	}
}
