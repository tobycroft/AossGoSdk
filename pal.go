package AossGoSdk

import (
	"errors"

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

func (self PalWorld) ShowPlayers() (users []OnlineUser, err error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/rcon/palworld/players", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat onlineUser
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		users = dat.Data
		return
	} else {
		return nil, errors.New(resp.Echo)
	}
}

func (self PalWorld) Kick(id any) error {
	ts := time.Now().Unix()
	param := map[string]any{
		"id":   id,
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/rcon/palworld/kick", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return err
	}
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}

func (self PalWorld) Ban(id any) error {
	ts := time.Now().Unix()
	param := map[string]any{
		"id":   id,
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/rcon/palworld/ban", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return err
	}
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}

func (self PalWorld) Ping() error {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	rets := new(Net.Post).PostFormDataAny(baseUrls+"/v1/rcon/palworld/ping", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return err
	}
	if rs.Code == 0 {
		return nil
	} else {
		return errors.New(rs.Echo)
	}
}
