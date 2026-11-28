package AossGoSdk

import (
	"errors"
	"time"

	"github.com/tobycroft/Calc"
)

type File struct {
	Token     string
	RemoteUrl string
	Ip        string
}

func (self *File) SetIp(ip string) *File {
	self.Ip = ip
	return self
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
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
	Url       string `json:"url"`
	UploadUrl string `json:"upload_url"`
}

func (self *File) GetUploadToken() (FileData, error) {
	ts := time.Now().Unix()
	sign := Calc.Md5(self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"token":     self.Token,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
	}
	if self.Ip != "" {
		param["ip"] = self.Ip
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := aossPostFormAny(remoteUrl+"/v2/file/token/create", nil, param, nil, nil)
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
	sign := Calc.Md5(self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"token":     self.Token,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
	}
	if self.Ip != "" {
		param["ip"] = self.Ip
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := aossPostFormAny(remoteUrl+"/v2/file/token/upload_url", nil, param, nil, nil)
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

type FileHashData struct {
	Src         string  `json:"src"`
	Url         string  `json:"url"`
	Surl        string  `json:"surl"`
	Name        string  `json:"name"`
	Mime        string  `json:"mime"`
	Path        string  `json:"path"`
	Ext         string  `json:"ext"`
	Size        int64   `json:"size"`
	Md5         string  `json:"md5"`
	Sha1        string  `json:"sha1"`
	Width       int64   `json:"width"`
	Height      int64   `json:"height"`
	Duration    float64 `json:"duration"`
	DurationStr string  `json:"duration_str"`
	Bitrate     float64 `json:"bitrate"`
	// V3 分发地址：默认与 Url 同为 attachment 永久地址；带 ttl 查询时为 file 临时地址
	FileUrl       string `json:"file_url"`
	FileExpiredAt int64  `json:"file_expired_at"`
}

type fileHashRet struct {
	Code int64        `json:"code"`
	Data FileHashData `json:"data"`
	Echo string       `json:"echo"`
}

func (self *File) GetUploadHashUrl() (FileUrlData, error) {
	ts := time.Now().Unix()
	sign := Calc.Md5(self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"token":     self.Token,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
	}
	if self.Ip != "" {
		param["ip"] = self.Ip
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := aossPostFormAny(remoteUrl+"/v2/file/token/upload_url_hash", nil, param, nil, nil)
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

func (self *File) GetUploadedFileUrlByHash(hash string) (FileHashData, error) {
	ts := time.Now().Unix()
	sign := Calc.Md5(self.Token + Calc.Any2String(ts))

	param := map[string]any{
		"token":     self.Token,
		"timestamp": Calc.Any2String(ts),
		"sign":      sign,
		"hash":      hash,
	}

	remoteUrl := baseUrls
	if self.RemoteUrl != "" {
		remoteUrl = self.RemoteUrl
	}

	rets := aossPostFormAny(remoteUrl+"/v2/file/token/hash_query", nil, param, nil, nil)
	var rs ret_std
	err := rets.RetJson(&rs)
	if err != nil {
		return FileHashData{}, err
	}
	if rs.Code == 0 {
		var dat fileHashRet
		err = rets.RetJson(&dat)
		if err != nil {
			return FileHashData{}, err
		}
		return dat.Data, nil
	}
	return FileHashData{}, errors.New(rs.Echo)
}
