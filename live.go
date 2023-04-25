package AossGoSdk

import (
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

type liveStructCreateAll struct {
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
func (self *Lcic) CreateAll(title any) (LiveStructCreateAll, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"title": title,
		"ts":    ts,
		"name":  self.Name,
		"sign":  Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/live/room/create_all", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LiveStructCreateAll{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return LiveStructCreateAll{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat liveStructCreateAll
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LiveStructCreateAll{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LiveStructCreateAll{}, errors.New(resp.Echo)
		}
	}
}

type liveStructCreateRoom struct {
	Code int                  `json:"code"`
	Data LiveStructCreateRoom `json:"data"`
	Echo string               `json:"echo"`
}

type LiveStructCreateRoom struct {
	Rtmp        string `json:"rtmp"`
	Domain      string `json:"domain"`
	ObsServer   string `json:"obs_server"`
	StreamCode  string `json:"stream_code"`
	Webrtc      string `json:"webrtc"`
	Srt         string `json:"srt"`
	RtmpOverSrt string `json:"rtmp_over_srt"`
}

// CreateRoom创建房间并返回房间数据
func (self *Lcic) CreateRoom(title any) (LiveStructCreateRoom, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"title": title,
		"ts":    ts,
		"name":  self.Name,
		"sign":  Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/live/room/create", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LiveStructCreateRoom{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return LiveStructCreateRoom{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat liveStructCreateRoom
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LiveStructCreateRoom{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LiveStructCreateRoom{}, errors.New(resp.Echo)
		}
	}
}

type liveStructPlayUrl struct {
	Code int               `json:"code"`
	Data LiveStructPlayUrl `json:"data"`
	Echo string            `json:"echo"`
}

type LiveStructPlayUrl struct {
	PlayDomain string `json:"play_domain"`
	PlayFlv    string `json:"play_flv"`
	PlayHls    string `json:"play_hls"`
	PlayRtmp   string `json:"play_rtmp"`
}

// CreateRoom创建房间并返回房间数据
func (self *Lcic) GetPlayUrl() (LiveStructPlayUrl, error) {
	ts := time.Now().Unix()
	param := map[string]any{
		"ts":   ts,
		"name": self.Name,
		"sign": Calc.Md5(self.Token + Calc.Any2String(ts)),
	}
	ret, err := Net.Post(baseUrls+"/v1/live/room/play_url", map[string]interface{}{
		"token": self.Name,
	}, param, nil, nil)
	//fmt.Println(ret, err)
	if err != nil {
		return LiveStructPlayUrl{}, err
	} else {
		var resp ret_std
		err = jsoniter.UnmarshalFromString(ret, &resp)
		if err != nil {
			return LiveStructPlayUrl{}, errors.New(ret)
		}
		if resp.Code == 0 {
			var dat liveStructPlayUrl
			err = jsoniter.UnmarshalFromString(ret, &dat)
			if err != nil {
				return LiveStructPlayUrl{}, errors.New(ret)
			}
			return dat.Data, nil
		} else {
			return LiveStructPlayUrl{}, errors.New(resp.Echo)
		}
	}
}
