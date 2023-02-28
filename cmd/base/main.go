package main

import (
	"net"
	"simple_tiktok/cmd/base/dal"
	"simple_tiktok/kitex_gen/base/baseservice"
	"simple_tiktok/util/mw"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"simple_tiktok/util/consts"
)

func Init() {
	//调用 dal.Init() 进行数据访问层的初始化。
	//然后使用 klogrus.NewLogger() 进行日志的初始化，
	//将其设置为 klog 的日志实现，并设置 klog 的日志级别为 klog.LevelInfo。
	dal.Init()
	// klog init
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
	addr, err := net.ResolveTCPAddr(consts.TCP, consts.BaseServiceAddr)
	if err != nil {
		panic(err)
	}
	Init()
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.BaseServiceName),
		provider.WithExportEndpoint(consts.ExportEndpoint),
		provider.WithInsecure(),
	)
	//在创建 userservice.Server 实例时，通过调用不同的函数来设置服务器的一些参数，
	//比如设置服务器的地址、注册中心、限流、中间件、追踪器等。
	//最后调用 svr.Run() 启动服务器
	svr := baseservice.NewServer(new(BaseServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithMiddleware(mw.CommonMiddleware),
		server.WithMiddleware(mw.ServerMiddleware),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.BaseServiceName}),
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
