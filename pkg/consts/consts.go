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

	BaseServiceAddr     = ":9000"
	InteractServiceAddr = ":9001"
	SocialServiceAddr   = ":9002"

	ExportEndpoint = ":4317"
	ETCDAddress    = "127.0.0.1:2379"

	ExchangeName = "comments_exchange"
)
