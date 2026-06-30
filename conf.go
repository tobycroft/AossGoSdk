package AossGoSdk

var baseUrl = "http://upload.tuuz.cc:81"
var baseUrls = "https://upload.tuuz.cc:433"

func Wechat_conf_set(BaseUrl, BaseUrls, CDNUrl, CDNUrls string) {
	baseUrl = BaseUrl
	baseUrls = BaseUrls
}
