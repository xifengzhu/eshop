package wechat

import (
	"net/url"
)

func CodeToSession(code string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("appid", wechatAppId)
	params.Set("secret", wechatSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	path := "https://api.weixin.qq.com/sns/jscode2session"
	return httpGet(path, params)
}
