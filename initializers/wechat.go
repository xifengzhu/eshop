package initializers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/xifengzhu/eshop/initializers/setting"

	"github.com/gomodule/redigo/redis"
)

var wechatAppId string
var wechatSecret string
var redisConn redis.Conn

func init() {
	redisConn, _ = redis.DialURL(setting.RedisUrl)
	wechatAppId = setting.WechatAppId
	wechatSecret = setting.WechatSecret
}

func CodeToSession(code string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("appid", wechatAppId)
	params.Set("secret", wechatSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	path := "https://api.weixin.qq.com/sns/jscode2session"
	return httpGet(path, params)
}

func GetWxAccessToken() (token string, err error) {
	token = getAccessTokenFromRedis()
	if token != "" {
		return token, nil
	}
	token, err = RefreshWxAccessToken()
	return
}

func RefreshWxAccessToken() (string, error) {
	var token string
	params := url.Values{}
	params.Set("appid", wechatAppId)
	params.Set("secret", wechatSecret)
	params.Set("grant_type", "client_credential")

	path := "https://api.weixin.qq.com/cgi-bin/token"
	result, err := httpGet(path, params)

	if err != nil {
		return "", err
	}

	token = result["access_token"].(string)
	_, err = redisConn.Do("SET", "eshop:wechat_access_token", token)
	if err != nil {
		log.Println("set eshop:wechat_access_token err:", err)
	}
	_, err = redisConn.Do("EXPIRE", "eshop:wechat_access_token", 6000)
	if err != nil {
		log.Println("set expire wechat_access_token err:", err)
	}
	return token, nil
}

func getAccessTokenFromRedis() (token string) {
	reply, err := redis.Values(redisConn.Do("MGET", "eshop:wechat_access_token"))
	if err != nil {
		log.Println("get wechat_access_token err:", err)
	}
	if _, err := redis.Scan(reply, &token); err != nil {
		log.Println("scan reply err:", err)
	}
	log.Println("access token:", token)
	return token
}

func httpGet(path string, params url.Values) (map[string]interface{}, error) {
	var result map[string]interface{}
	Url, err := url.Parse(path)
	if err != nil {
		panic(err.Error())
	}

	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	log.Println("urlPath:", urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	} else {
		log.Println("http resp", resp)
		defer resp.Body.Close()
		s, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(s), &result)
		log.Println("request result", result)
		return result, nil
	}
}
