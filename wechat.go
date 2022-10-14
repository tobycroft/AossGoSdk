package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

type wechatRet struct {
	Code int64
	Echo string
}

type wechatWxaUnlimitedFile struct {
	Code int64
	Data string
	Echo string
}

func Wechat_wxa_unlimited_file(project, data, page string) (string, error) {
	post := map[string]any{
		"data": data,
		"page": page,
	}
	ret, err := Net.Post(baseUrl+unlimited_file, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return "", err
	}
	var wwuf wechatWxaUnlimitedFile
	err = jsoniter.UnmarshalFromString(ret, &wwuf)
	if err != nil {
		return "", err
	}
	if wwuf.Code == 0 {
		return wwuf.Data, nil
	} else {
		return wwuf.Data, errors.New(wwuf.Echo)
	}
}

func Wechat_wxa_unlimited_raw(project, data, page string) ([]byte, error) {
	post := map[string]any{
		"data": data,
		"page": page,
	}
	ret, err := Net.Post(baseUrl+unlimited_base64, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return nil, err
	}
	var wwuf wechatWxaUnlimitedFile
	err = jsoniter.UnmarshalFromString(ret, &wwuf)
	if err != nil {
		return nil, err
	}
	if wwuf.Code == 0 {
		return decode(wwuf.Data)
	} else {
		return nil, errors.New(wwuf.Echo)
	}
}

type wechatWxaRet struct {
	Data WechatWxaScene
}

type WechatWxaScene struct {
	Key  string
	Val  string
	Page string
	Path string
	Url  string
}

func Wechat_wxa_scene(project, scene string) (WechatWxaScene, error) {
	post := map[string]any{
		"scene": scene,
	}
	ret, err := Net.Post(baseUrl+wxa_scene, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return WechatWxaScene{}, err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return WechatWxaScene{}, errors.New(ret)
	}
	if resp.Code == 0 {
		var dat wechatWxaRet
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return WechatWxaScene{}, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return WechatWxaScene{}, errors.New(resp.Echo)
	}
}

type WechatSnsJscode2sessionRet struct {
	Data WechatSnsJscode2session
}

type WechatSnsJscode2session struct {
	SessionKey string
	Unionid    string
	Openid     string
}

func Wechat_sns_jscode2session(project, js_code string) (WechatSnsJscode2session, error) {
	post := map[string]any{
		"js_code": js_code,
	}
	ret, err := Net.Post(baseUrl+jscode, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return WechatSnsJscode2session{}, err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return WechatSnsJscode2session{}, errors.New(ret)
	}
	if resp.Code == 0 {
		var dat WechatSnsJscode2sessionRet
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return WechatSnsJscode2session{}, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return WechatSnsJscode2session{}, errors.New(resp.Echo)
	}
}

type WechatWxaGEtUserPhoneNumberRet struct {
	Data WechatWxaGEtUserPhoneNumber
}

type WechatWxaGEtUserPhoneNumber struct {
	PhoneNumber     string      `json:"phoneNumber"`
	PurePhoneNumber string      `json:"purePhoneNumber"`
	CountryCode     string      `json:"countryCode"`
	Watermark       interface{} `json:"watermark"`
}

func Wechat_wxa_getuserphonenumber(project, code string) (WechatWxaGEtUserPhoneNumber, error) {
	post := map[string]any{
		"code": code,
	}
	ret, err := Net.Post(baseUrl+getuserphonenumber, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return WechatWxaGEtUserPhoneNumber{}, err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return WechatWxaGEtUserPhoneNumber{}, errors.New(ret)
	}
	if resp.Code == 0 {
		var dat WechatWxaGEtUserPhoneNumberRet
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return WechatWxaGEtUserPhoneNumber{}, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return WechatWxaGEtUserPhoneNumber{}, errors.New(resp.Echo)
	}
}
