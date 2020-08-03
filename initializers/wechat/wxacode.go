package wechat

import (
	"encoding/json"
	"net/http"
)

func GetWxaCodeUnLimit(page string, scene string, width int, is_hyaline bool) (*http.Response, error) {

	params := make(map[string]interface{})
	params["page"] = page
	params["scene"] = scene
	params["width"] = width
	params["is_hyaline"] = is_hyaline
	jsonParams, _ := json.Marshal(params)

	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + getAccessTokenFromRedis()
	return httpPost(url, jsonParams)
}
