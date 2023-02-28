package consts

const (
	MySQLDefaultDSN = "jiyeon:1234@tcp(localhost:3307)/tiktok?charset=utf8&parseTime=True&loc=Local"
	UserTableName   = "users"
	VideoTableName  = "Videos"

	ApiServiceName = "api"

	BaseServiceName     = "base"
	InteractServiceName = "interact"

	TCP = "tcp"

	BaseServiceAddr     = ":9000"
	InteractServiceAddr = ":9001"

	ExportEndpoint = ":4317"
	ETCDAddress    = "127.0.0.1:2379"
	DefaultLimit   = 10
)
