package initializers

import (
	"github.com/gocraft/work"
	"github.com/xifengzhu/eshop/initializers/setting"
)

var JobEnqueuer *work.Enqueuer

func init() {
	JobEnqueuer = work.NewEnqueuer(setting.RedisNamespace, RedisPool)
}
