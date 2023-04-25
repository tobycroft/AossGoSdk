package AossGoSdk

import (
	"encoding/json"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
	"time"
)

type Live struct {
	Token string
	Code  string
}

type liveStructCreateRoom struct {
	Code int                 `json:"code"`
	Data LiveStructCreateAll `json:"data"`
	Echo string              `json:"echo"`
}

type LiveStructCreateAll struct {
	Rtmp        string `json:"rtmp"`
	Domain      string `json:"domain"`
	PlayDomain  string `json:"play_domain"`
	ObsServer   string `json:"obs_server"`
	StreamCode  string `json:"stream_code"`
	Webrtc      string `json:"webrtc"`
	Srt         string `json:"srt"`
	RtmpOverSrt string `json:"rtmp_over_srt"`
	PlayFlv     string `json:"play_flv"`
	PlayHls     string `json:"play_hls"`
	PlayRtmp    string `json:"play_rtmp"`
}

// RoomCreate TeacherId:老师ID | StartTime:开始时间int | EndTime:结束时间int | Name:房间名称
func (self *Lcic) CreateAll() (LcicStructCreateRoom, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/lcic/room/create", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LcicStructCreateRoom{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return LcicStructCreateRoom{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat liveStructCreateRoom
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LcicStructCreateRoom{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LcicStructCreateRoom{}, errors.New(resp.Echo)
		}
	}
}

// Lcic_RoomModify RoomId:房间ID | TeacherId:老师ID | StartTime:开始时间int | EndTime:结束时间int | Name:房间名称
func (self *Lcic) Lcic_RoomModify(RoomId, TeacherId, StartTime, EndTime, Name interface{}) (bool, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"RoomId":    RoomId,
		"Name":      Name,
		"TeacherId": TeacherId,
		"StartTime": StartTime,
		"EndTime":   EndTime,
		"ts":        ts,
		"name":      self.Name,
		"sign":      Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/lcic/room/create", map[string]interface{}{
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
	ret, err := Net.Post(baseUrls+"/v1/lcic/room/delete", map[string]interface{}{
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
	ret, err := Net.Post(baseUrls+"/v1/lcic/room/delete", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LcicStructLinkUrl{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
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
