package AossGoSdk

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	Net "github.com/tobycroft/TuuzNet"
)

func Wechat_get_accessToken(project string) error {
	ret, err := Net.Post(baseUrl+get_accesstoken, map[string]interface{}{
		"token": project,
	}, nil, nil, nil)
	if err != nil {
		return err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return errors.New(ret)
	}
	if resp.Code == 0 {
		var dat wechatDataMapStringInterface
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return errors.New(ret)
		}
		postData := dat.Data["postdata"].(map[string]interface{})
		ret, err = Net.Post(dat.Data["address"].(string), nil, postData, nil, nil)
		if err != nil {
			return err
		}
		ret, err = Net.Post(baseUrl+set_accesstoken, map[string]interface{}{
			"token": project,
		}, map[string]interface{}{
			"data": ret,
		}, nil, nil)
		if err != nil {
			return err
		}
		var resp2 wechatRet
		err = jsoniter.UnmarshalFromString(ret, &resp2)
		if err != nil {
			return errors.New(ret)
		}
		if resp2.Code == 0 {
			return nil
		} else {
			return errors.New(resp2.Echo)
		}
	} else {
		return errors.New(resp.Echo)
	}
}

type wechatDataSlices struct {
	Data []string
}

type wechatDataMapStringInterface struct {
	Data map[string]interface{}
}

func Wechat_offi_get_user_list(project string) ([]string, error) {
	ret, err := Net.Post(baseUrl+offi_user_list, map[string]interface{}{
		"token": project,
	}, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return nil, errors.New(ret)
	}
	if resp.Code == 0 {
		var dat wechatDataSlices
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return nil, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return nil, errors.New(resp.Echo)
	}
}

type wechatDataUserInfo struct {
	Data WechatUserInfo
}

type WechatUserInfo struct {
	subscribe      int64
	openid         string
	nickname       string
	sex            int64
	headimgurl     string
	subscribe_time int64
}

func Wechat_offi_get_user_info(project, openid string) (WechatUserInfo, error) {
	post := map[string]any{
		"openid": openid,
	}
	ret, err := Net.Post(baseUrl+offi_user_list, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return WechatUserInfo{}, err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return WechatUserInfo{}, errors.New(ret)
	}
	if resp.Code == 0 {
		var dat wechatDataUserInfo
		err = jsoniter.UnmarshalFromString(ret, &dat)
		if err != nil {
			return WechatUserInfo{}, errors.New(ret)
		}
		return dat.Data, nil
	} else {
		return WechatUserInfo{}, errors.New(resp.Echo)
	}
}

// Wechat_offi_openidUrl:redirect_url无需urlencode，AOSS会自动urlencode
func Wechat_offi_openidUrl(project, redirect_uri, response_type, scope, state string, show_in_qrcode bool) (string, error) {
	post := map[string]any{
		"redirect_uri":  redirect_uri,
		"response_type": response_type,
		"scope":         scope,
		"state":         state,
		"png":           show_in_qrcode,
	}
	ret, err := Net.Post(baseUrl+offi_openid_url, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return "", err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return "", errors.New(ret)
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = jsoniter.UnmarshalFromString(ret, &wwuf)
		if err != nil {
			return "", err
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}

// Wechat_offi_openid_from_code:通过微信前置授权取回code后，通过该code获取openid
func Wechat_offi_openid_from_code(project, code any) (string, error) {
	post := map[string]any{
		"code": code,
	}
	ret, err := Net.Post(baseUrl+offi_openid_url, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	if err != nil {
		return "", err
	}
	var resp wechatRet
	err = jsoniter.UnmarshalFromString(ret, &resp)
	if err != nil {
		return "", errors.New(ret)
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = jsoniter.UnmarshalFromString(ret, &wwuf)
		if err != nil {
			return "", err
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}
