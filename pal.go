package AossGoSdk

import (
	"encoding/json"
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

func (self PalWorld) ShowPlayers() ([]OnlineUser, error) {
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

func (self PalWorld) Kick(id any) error {
	ts := time.Now().Unix()
	param := map[string]any{
		"id":   id,
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/rcon/palworld/kick", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return err
	} else {
		var rs ret_std
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

func (self PalWorld) Ban(id any) error {
	ts := time.Now().Unix()
	param := map[string]any{
		"id":   id,
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/rcon/palworld/ban", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return err
	} else {
		var rs ret_std
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

func (self PalWorld) Ping() error {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/rcon/palworld/ping", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return err
	} else {
		var rs ret_std
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
