package utils

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mojocn/base64Captcha"
	"github.com/xifengzhu/eshop/helpers/setting"
	"log"
	"strings"
)

type redisStore struct {
	redisClient redis.Conn
}

// redisStore Implement the store interface(Set, Get, Verify)
func (s *redisStore) Set(id string, value string) {
	var err error
	_, err = s.redisClient.Do("SETEX", id, 600, value)
	if err != nil {
		panic(err)
	}
}

// redisStore implementing Get method of  Store interface
func (s *redisStore) Get(id string, clear bool) (value string) {
	reply, err := redis.Values(redisConn.Do("MGET", id))
	if err != nil {
		log.Println("====get %s err=====", id)
	}
	if _, err := redis.Scan(reply, &value); err != nil {
		log.Println("scan reply err:", err)
	}
	return
}

func (s *redisStore) Verify(id, answer string, clear bool) (match bool) {
	vv := s.Get(id, clear)
	vv = strings.TrimSpace(vv)
	return vv == strings.TrimSpace(answer)
}

var store *redisStore

func init() {
	// init redis store
	redisConn, _ = redis.DialURL(setting.RedisUrl)
	store = &redisStore{redisConn}

	// SetCustomStore is not working
	base64Captcha.DefaultMemStore = store
}

func CaptchaGenerate() (id, b64s string, err error) {
	driver := base64Captcha.DefaultDriverDigit
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err = c.Generate()
	return
}

func CaptchaVerify(id, value string) bool {
	if store.Verify(id, value, true) {
		return true
	}
	return false
}
