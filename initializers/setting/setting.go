package setting

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

var (
	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize        int
	JwtSecret       string
	Domain          string
	RuntimeRootPath string

	WechatAppId     string
	WechatSecret    string
	MchID           string
	MchKey          string
	PayNotifyUrl    string
	RefundNotifyUrl string

	LogFilePath string
	LogFileName string

	DatabaseType        string
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	DatabaseHost        string
	DatabaseTablePrefix string
	DatabasePort        string

	RedisUrl       string
	RedisNamespace string

	QiniuAccessKey    string
	QiniuSecretKey    string
	QiniuBucket       string
	QiniuBucketDomain string

	ExportSavePath string
)

func init() {
	env := os.Getenv("APP_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env

	LoadBase()
	LoadServer()
	LoadApp()
	LoadWechat()
	LoadLog()
	LoadDatabase()
	LoadRedis()
	LoadQiniu()
	LoadExport()
}

func LoadBase() {
	RunMode = os.Getenv("RUN_MODE")
}

func LoadServer() {
	HTTPPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	rtimeout, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT"))
	wtimeout, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT"))
	ReadTimeout = time.Duration(rtimeout) * time.Second
	WriteTimeout = time.Duration(wtimeout) * time.Second
}

func LoadApp() {
	JwtSecret = os.Getenv("JWT_SECRET")
	PageSize, _ = strconv.Atoi(os.Getenv("PAGE_SIZE"))
	Domain = os.Getenv("DOMAIN")
	RuntimeRootPath = os.Getenv("RUNTIME_ROOT_PATH")
}

func LoadWechat() {
	WechatAppId = os.Getenv("WECHAT_APPID")
	WechatSecret = os.Getenv("WECHAT_SECRET")

	MchID = os.Getenv("MCH_ID")
	MchKey = os.Getenv("MCH_KEY")
	PayNotifyUrl = os.Getenv("PAY_NOTIFY_URL")
	RefundNotifyUrl = os.Getenv("REFUND_NOTIFY_URL")
}

func LoadLog() {
	LogFilePath = os.Getenv("LOG_FILE_PATH")
	LogFileName = os.Getenv("LOG_FILENAME")
}

func LoadDatabase() {
	DatabaseType = os.Getenv("TYPE")
	DatabaseName = os.Getenv("NAME")
	DatabaseUser = os.Getenv("USER")
	DatabasePassword = os.Getenv("PASSWORD")
	DatabaseHost = os.Getenv("HOST")
	DatabaseTablePrefix = os.Getenv("TABLE_PREFIX")
	DatabasePort = os.Getenv("PORT")
}

func LoadRedis() {
	RedisUrl = os.Getenv("REDIS_URL")
	RedisNamespace = os.Getenv("REDIS_NAMESPACE")
}

func LoadQiniu() {
	QiniuAccessKey = os.Getenv("ACCESS_KEY")
	QiniuSecretKey = os.Getenv("SECRET_KEY")
	QiniuBucket = os.Getenv("BUCKET")
	QiniuBucketDomain = os.Getenv("BUCKET_DOMAIN")
}

func LoadExport() {
	ExportSavePath = os.Getenv("EXPORT_SAVE_PATH")
}
