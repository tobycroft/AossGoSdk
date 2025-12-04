package AossGoSdk

import (
	"errors"

	"github.com/bytedance/sonic"
	Net "github.com/tobycroft/TuuzNet"
)

// Wechat_get_accessToken:获取AccessToken，不推荐，建议使用0.0.0.0/0设定公众号
//func Wechat_get_accessToken(project string) error {
//	rets := new(Net.Post).PostFormDataAny(baseUrl+get_accesstoken, map[string]interface{}{
//		"token": project,
//	}, nil, nil, nil)
//	if err != nil {
//		return err
//	}
//	var resp ret_std
//	err = rets.RetJson(&resp)
//	if err != nil {
//		return errors.New(ret)
//	}
//	if resp.Code == 0 {
//		var dat wechatDataMapStringInterface
//		err = jsoniter.UnmarshalFromString(ret, &dat)
//		if err != nil {
//			return errors.New(ret)
//		}
//		postData := dat.Data["postdata"].(map[string]interface{})
//		ret, err = new(Net.Post).PostFormDataAny(dat.Data["address"].(string), nil, postData, nil, nil)
//		if err != nil {
//			return err
//		}
//		ret, err = new(Net.Post).PostFormDataAny(baseUrl+set_accesstoken, map[string]interface{}{
//			"token": project,
//		}, map[string]interface{}{
//			"data": ret,
//		}, nil, nil)
//		if err != nil {
//			return err
//		}
//		var resp2 ret_std
//		err = jsoniter.UnmarshalFromString(ret, &resp2)
//		if err != nil {
//			return errors.New(ret)
//		}
//		if resp2.Code == 0 {
//			return nil
//		} else {
//			return errors.New(resp2.Echo)
//		}
//	} else {
//		return errors.New(resp.Echo)
//	}
//}

type wechatDataSlices struct {
	Data []string
}

type wechatDataMapStringInterface struct {
	Data map[string]interface{}
}

// Wechat_offi_get_user_list:获取已关注用户的openid列表，slices仅限openid无法区别谁是谁
func Wechat_offi_get_user_list(project string) (str []string, err error) {
	rets := new(Net.Post).PostFormDataAny(baseUrl+offi_user_list, map[string]interface{}{
		"token": project,
	}, nil, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatDataSlices
		err = rets.RetJson(&dat)
		if err != nil {
			return
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

// Wechat_offi_get_user_info:微信获取某个用户的详细情况（但是有可能获取不到）
func Wechat_offi_get_user_info(project, openid string) (info WechatUserInfo, err error) {
	post := map[string]any{
		"openid": openid,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+offi_user_list, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var dat wechatDataUserInfo
		err = rets.RetJson(&dat)
		if err != nil {
			return
		}
		return dat.Data, nil
	} else {
		return WechatUserInfo{}, errors.New(resp.Echo)
	}
}

// Wechat_offi_openidUrl:redirect_url无需urlencode，AOSS会自动urlencode
func Wechat_offi_openidUrl(project, redirect_uri, response_type, scope, state string, show_in_qrcode bool) (str string, err error) {
	post := map[string]any{
		"redirect_uri":  redirect_uri,
		"response_type": response_type,
		"scope":         scope,
		"state":         state,
		"png":           show_in_qrcode,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+offi_openid_url, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
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

// Wechat_offi_openid_from_code:通过微信前置授权取回code后，通过该code获取openid
func Wechat_offi_openid_from_code(project, code any) (str string, err error) {
	post := map[string]any{
		"code": code,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+offi_openid_aquire_from_code, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = rets.RetJson(&wwuf)
		if err != nil {
			return
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}

type Wechat_template_data_struct struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// Wechat_template_send:微信发送模版功能
func Wechat_template_send(project, openid, template_id, url interface{}, data map[string]Wechat_template_data_struct) (str string, err error) {
	dat, _ := sonic.MarshalString(data)
	post := map[string]any{
		"openid":      openid,
		"template_id": template_id,
		"url":         url,
		"data":        dat,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+offi_template_send, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
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

// Wechat_snsAuth:验证前端获取到的accesstoken和openid是否正确匹配
func Wechat_snsAuth(project, access_token, openid interface{}) (str string, err error) {
	post := map[string]any{
		"openid":       openid,
		"access_token": access_token,
	}
	rets := new(Net.Post).PostFormDataAny(baseUrl+sns_auth, map[string]interface{}{
		"token": project,
	}, post, nil, nil)
	var resp ret_std
	err = rets.RetJson(&resp)
	if err != nil {
		return
	}
	if resp.Code == 0 {
		var wwuf wechatStringData
		err = rets.RetJson(&wwuf)
		if err != nil {
			return
		}
		return wwuf.Data, nil
	} else {
		return "", errors.New(resp.Echo)
	}
}
