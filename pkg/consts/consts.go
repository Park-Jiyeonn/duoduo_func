package consts

const (
	MySQLDefaultDSN = "jiyeon:1234@tcp(localhost:3307)/tiktok?charset=utf8&parseTime=True&loc=Local"
	RabbitMQDSN     = "amqp://admin:password@localhost:5672/"

	UserTableName  = "users"
	VideoTableName = "Videos"

	ApiServiceName = "api"

	BaseServiceName     = "base"
	InteractServiceName = "interact"
	SocialServiceName   = "social"

	TCP = "tcp"

	BaseServiceAddr     = ":10086"
	InteractServiceAddr = ":10087"
	SocialServiceAddr   = ":10088"

	ExportEndpoint = ":4317"
	ETCDAddress    = "localhost:2379"

	ExchangeName = "comments_exchange"

	Limits_per_sec = 20
)
