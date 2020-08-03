package wechat

import (
	"bytes"
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

func httpPost(url string, params []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp, err
}
