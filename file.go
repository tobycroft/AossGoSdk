package AossGoSdk

import (
	"errors"
	"time"

	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
)

type File struct {
	Appid     string
	Token     string
	RemoteUrl string
}

type fileCreateRet struct {
	Code int64    `json:"code"`
	Data FileData `json:"data"`
	Echo string   `json:"echo"`
}

type FileData struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

type fileUrlRet struct {
	Code int64       `json:"code"`
	Data FileUrlData `json:"data"`
	Echo string      `json:"echo"`
}

type FileUrlData struct {
	UploadUrl string `json:"upload_url"`
}

func (self *File) GetUploadToken() (FileData, error) {
	ts := time.Now().Unix()
	sign := Calc.Md5(self.Appid + self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"appid":     self.Appid,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := new(Net.Post).PostFormDataAny(remoteUrl+"/v2/file/token/create", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return FileData{}, err
	}
	if rs.Code == 0 {
		var dat fileCreateRet
		err = rets.RetJson(&dat)
		if err != nil {
			return FileData{}, err
		}
		return dat.Data, nil
	}
	return FileData{}, errors.New(rs.Echo)
}

func (self *File) GetUploadUrl() (FileUrlData, error) {
	ts := time.Now().Unix()
	sign := Calc.Md5(self.Appid + self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"appid":     self.Appid,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := new(Net.Post).PostFormDataAny(remoteUrl+"/v2/file/token/upload_url", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return FileUrlData{}, err
	}
	if rs.Code == 0 {
		var dat fileUrlRet
		err = rets.RetJson(&dat)
		if err != nil {
			return FileUrlData{}, err
		}
		return dat.Data, nil
	}
	return FileUrlData{}, errors.New(rs.Echo)
}