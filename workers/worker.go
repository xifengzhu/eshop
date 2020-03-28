package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/xifengzhu/eshop/helpers/setting"
	"log"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.DialURL(setting.RedisUrl)
	},
}

type Context struct {
	customerID int64
}

func init() {
	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, setting.RedisNamespace, redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)

	// Map the name of jobs to handler functions
	pool.Job("send_email", (*Context).SendEmail)

	// Customize options:
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)

	pool.PeriodicallyEnqueue("* 1 * * * *", "refresh_wechat_access_token")
	pool.Job("refresh_wechat_access_token", (*Context).RefreshWechatAccessToken)
	// Start processing jobs
	pool.Start()

}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) FindCustomer(job *work.Job, next work.NextMiddlewareFunc) error {
	// If there's a customer_id param, set it in the context for future middleware and handlers to use.
	if _, ok := job.Args["customer_id"]; ok {
		c.customerID = job.ArgInt64("customer_id")
		if err := job.ArgError(); err != nil {
			return err
		}
	}

	return next()
}
