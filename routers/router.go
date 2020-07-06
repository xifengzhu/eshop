package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/initializers/setting"
	adminApi "github.com/xifengzhu/eshop/routers/admin_api"
	appApi "github.com/xifengzhu/eshop/routers/app_api"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/xifengzhu/eshop/docs"

	log "github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
	"strings"
	// "github.com/xifengzhu/eshop/middleware/logger"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func initSentry() {
	raven.SetDSN("https://c057910274d2496c90f7ab012072310f:5a6f7f0946924fc4a8b5fccedc92507c@sentry.beansmile-dev.com/55")
	raven.SetEnvironment(setting.RunMode)
}

func init() {
	initSentry()
	initLog()
}

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(gin.Recovery())

	if isProduction() {
		r.Use(globalRecover)
	}

	gin.SetMode(setting.RunMode)

	r.StaticFS("/static", http.Dir("./static"))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	appApi.InitAppAPI(r)
	adminApi.InitAdminAPI(r)
	return r
}

func initLog() {
	// Log as JSON instead of the default ASCII formatter.
	if isProduction() {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			DebugStack := ""
			for _, v := range strings.Split(string(debug.Stack()), "\n") {
				DebugStack += v + "<br/>"
			}
			body := c.Request.Method + "  " + c.Request.RequestURI + c.ClientIP() + "<br/>" + DebugStack

			raven.CaptureMessage(body, nil, nil)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器繁忙，请稍候再试！"})
		}
	}(c)
	c.Next()
}

func isProduction() bool {
	return setting.RunMode == "production"
}
