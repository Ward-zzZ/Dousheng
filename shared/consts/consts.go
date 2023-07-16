package consts

import "time"

const (
	IdentityKey = "id" //jwt携带的识别信息

	HlogFilePath = "./tmp/hlog/logs/" //日志的输出路径
	KlogFilePath = "./tmp/klog/logs/"

	ConfigPath = "./config.yaml"

	RedisExpireTime = time.Hour * 48
	MySqlDSN        = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	RabbitMqURI     = "amqp://%s:%s@%s:%d/" //消息队列

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0" // 本地和局域网都可以访问
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	TCP = "tcp"

	FreePortAddress = "localhost:0" // 0 means random port
	CorsAddress = "http://localhost:3000"

	ConsulCheckInterval                       = "7s"
	ConsulCheckTimeout                        = "5s"
	ConsulCheckDeregisterCriticalServiceAfter = "15s"

	// 不同服务的redis数据库
	RedisUserClientDB     = 0
	RedisVideoClientDB    = 0
	RedisRelationClientDB = 0
	RedisCommentClientDB  = 0
	RedisFavoriteClientDB = 0

	RedisFansThreshold = 1000

	UserSnowflakeNode  = 2
	VideoSnowflakeNode = 3

	SleepTime = time.Millisecond * 600 // 延时双删的时间

	Tiktok = "tiktok-demo" // sentinel

	VideoNumPerFeed = 10 // 每次获取feed的视频数量
)
