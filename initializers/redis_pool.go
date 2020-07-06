package initializers

import (
	"github.com/gomodule/redigo/redis"
	"github.com/xifengzhu/eshop/initializers/setting"
)

var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.DialURL(setting.RedisUrl)
	},
}
