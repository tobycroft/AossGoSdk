package AossGoSdk

import (
	"errors"

	Net "github.com/tobycroft/TuuzNet"
)

type ret_std struct {
	Code int64
	Echo string
}

type wechatStringData struct {
	Code int64
	Data string
	Echo string
}

// Wechat_wxa_unlimited_file:获取微信小程序二维码（302方法，推荐占用少）
func Wechat_wxa_unlimited_file(project, data, page string) (str string, err error) {
	post := map[string]any{
		"data": data,
		"page": page,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+unlimited_file, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return "", err
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = rets.RetJson(&wwuf)
		if err != nil {
			return "", err
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}

// Wechat_wxa_unlimited_raw:获取微信小程序二维码（文件流方法，不推荐会吃服务器内存）
func Wechat_wxa_unlimited_raw(project, data, page string) (b []byte, err error) {
	post := map[string]any{
		"data": data,
		"page": page,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+unlimited_base64, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return nil, err
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = rets.RetJson(&wwuf)
		if err != nil {
			return nil, err
		}
		return decode(wwuf.Data)
	} else {
		return nil, errors.New(resp.Echo)
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

// Wechat_wxa_scene:微信scene解析
func Wechat_wxa_scene(project, scene string) (sen WechatWxaScene, err error) {
	post := map[string]any{
		"scene": scene,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+wxa_scene, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatWxaRet
		err = rets.RetJson(&dat)
		return dat.Data, nil
	} else {
		return WechatWxaScene{}, errors.New(resp.Echo)
	}
}

type wechatSnsJscode2sessionRet struct {
	Data WechatSnsJscode2session
}

type WechatSnsJscode2session struct {
	SessionKey string
	Unionid    string
	Openid     string
}

// Wechat_sns_jscode2session:微信授权一键登录
func Wechat_sns_jscode2session(project, js_code string) (session WechatSnsJscode2session, err error) {
	post := map[string]any{
		"js_code": js_code,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+jscode, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatSnsJscode2sessionRet
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		return dat.Data, nil
	} else {
		return WechatSnsJscode2session{}, errors.New(resp.Echo)
	}
}

type wechatWxaGEtUserPhoneNumberRet struct {
	Data WechatWxaGEtUserPhoneNumber
}

type WechatWxaGEtUserPhoneNumber struct {
	PhoneNumber     string      `json:"phoneNumber"`
	PurePhoneNumber string      `json:"purePhoneNumber"`
	CountryCode     string      `json:"countryCode"`
	Watermark       interface{} `json:"watermark"`
}

// Wechat_wxa_getuserphonenumber:获取用户手机号
func Wechat_wxa_getuserphonenumber(project, code string) (num WechatWxaGEtUserPhoneNumber, err error) {
	post := map[string]any{
		"code": code,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+getuserphonenumber, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatWxaGEtUserPhoneNumberRet
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		return dat.Data, nil
	} else {
		return WechatWxaGEtUserPhoneNumber{}, errors.New(resp.Echo)
	}
}

type wechatWxaGenerateSchemeRet struct {
	Data WechatWxaGenerateScheme
}

type WechatWxaGenerateScheme struct {
	Openlink string
}

// Wechat_wxa_generatescheme:创建scheme-url地址
func Wechat_wxa_generatescheme(project, path, query string, is_expire bool, expire_interval int) (sche WechatWxaGenerateScheme, err error) {
	post := map[string]any{
		"path":            path,
		"query":           query,
		"is_expire":       is_expire,
		"expire_interval": expire_interval,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+generatescheme, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatWxaGenerateSchemeRet
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		return dat.Data, nil
	} else {
		return WechatWxaGenerateScheme{}, errors.New(resp.Echo)
	}
}
