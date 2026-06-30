package AossGoSdk

import (
	"errors"
	"time"

	"github.com/tobycroft/Calc"
	Net "github.com/tobycroft/TuuzNet"
)

type FileToken struct {
	Appid     string
	Token     string
	RemoteUrl string
}

type fileTokenCreateRet struct {
	Code int64         `json:"code"`
	Data FileTokenData `json:"data"`
	Echo string        `json:"echo"`
}

type FileTokenData struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

func (self *FileToken) GetUploadToken() (FileTokenData, error) {
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
		return FileTokenData{}, err
	}
	if rs.Code == 0 {
		var dat fileTokenCreateRet
		err = rets.RetJson(&dat)
		if err != nil {
			return FileTokenData{}, err
		}
		return dat.Data, nil
	}
	return FileTokenData{}, errors.New(rs.Echo)
}
