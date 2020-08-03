package initializers

import (
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
