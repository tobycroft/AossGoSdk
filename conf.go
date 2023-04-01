package AossGoSdk

var baseUrl = "http://upload.familyeducation.org.cn:81"
var baseUrls = "https://upload.familyeducation.org.cn:444"
var cdnUrl = "http://aoss.familyeducation.org.cn"
var cdnUrls = "https://aoss.familyeducation.org.cn"

func Wechat_conf_set(BaseUrl, BaseUrls, CDNUrl, CDNUrls string) {
	baseUrl = BaseUrl
	baseUrls = BaseUrls
	cdnUrl = CDNUrl
	cdnUrls = CDNUrls
}
