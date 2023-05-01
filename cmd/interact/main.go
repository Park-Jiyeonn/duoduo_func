package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"net"
	"simple_tiktok/dal"
	"simple_tiktok/kitex_gen/interact/interactservice"
	"simple_tiktok/pkg/consts"
	"simple_tiktok/util/mw"
)

func Init() {
	dal.Init()
	klog.SetLogger(kitexlogrus.NewLogger())
	klog.SetLevel(klog.LevelInfo)
}

func main() {
	//首先使用 etcd.NewEtcdRegistry() 函数创建一个 etcd 注册中心实例，
	//consts.ETCDAddress 是 etcd 服务器的地址。
	r, err := etcd.NewEtcdRegistry([]string{consts.ETCDAddress})
	if err != nil {
		panic(err)
	}
	//net.ResolveTCPAddr() 函数将onsts.UserServiceAddr解析为一个 TCP 地址类型。
	//接着调用 Init() 函数初始化一些配置，然后创建一个 userservice.Server 的实例
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.InteractServiceAddr)
	if err != nil {
		panic(err)
	}
	Init()
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.InteractServiceName),
		provider.WithExportEndpoint(consts.ExportEndpoint),
		provider.WithInsecure(),
	)
	//在创建 userservice.Server 实例时，通过调用不同的函数来设置服务器的一些参数，
	//比如设置服务器的地址、注册中心、限流、中间件、追踪器等。
	//最后调用 svr.Run() 启动服务器
	svr := interactservice.NewServer(new(InteractServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithMiddleware(mw.CommonMiddleware),
		server.WithMiddleware(mw.ServerMiddleware),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.InteractServiceName}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
