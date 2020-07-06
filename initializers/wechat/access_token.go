package wechat

import (
	"log"
	"net/url"

	"github.com/gomodule/redigo/redis"
)

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
