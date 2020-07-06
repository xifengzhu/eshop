package workers

import (
	"github.com/gocraft/work"
	// "github.com/xifengzhu/eshop/initializers/wechat"
	"log"
)

func (c *Context) RefreshWechatAccessToken(job *work.Job) error {
	log.Println("Cronjob: RefreshWechatAccessToken...")
	// wechat.RefreshWxAccessToken()
	return nil
}
