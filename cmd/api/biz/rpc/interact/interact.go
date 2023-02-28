package interact

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"simple_tiktok/kitex_gen/interact"
	"simple_tiktok/kitex_gen/interact/interactservice"
	"simple_tiktok/util/consts"
	"simple_tiktok/util/mw"
)

var interactClient interactservice.Client

func InitInteract() {
	//使用 etcd 作为服务发现和注册中心的解析器，将用户服务的地址从 etcd 中获取。
	r, err := etcd.NewEtcdResolver([]string{consts.ETCDAddress})
	if err != nil {
		panic(err)
	}

	//创建一个新的 OpenTelemetry 提供程序，用于生成跟踪和度量数据。
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(consts.ApiServiceName),
		provider.WithExportEndpoint(consts.ExportEndpoint),
		provider.WithInsecure(),
	)

	//创建一个新的客户端，该函数需要传递以下参数：
	//	用户服务的服务名称。
	//	etcd 解析器，用于解析服务的地址。
	//	最大连接数为1的 Mux 连接器，用于创建一个连接池并且仅使用一个连接。
	//	客户端中间件，用于处理请求和响应之前的拦截和操作。
	//	实例中间件，用于处理客户端的整个生命周期。
	//	追踪器中间件，用于收集调用的跟踪信息。
	//	客户端基本信息，包含服务名称等元数据。
	c, err := interactservice.NewClient(
		consts.InteractServiceName,
		client.WithResolver(r),
		client.WithMuxConnection(1),
		client.WithMiddleware(mw.CommonMiddleware), // 这个中间件的作用是用于打印一些 RPC 的信息，比如请求的参数、远程服务的名称和方法以及返回结果等。
		client.WithInstanceMW(mw.ClientMiddleware), // 这个中间件的作用是在客户端请求远程服务之前，打印出远程服务的地址、RPC超时时间和读写超时时间。这对于客户端调试和问题排查非常有用。
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.ApiServiceName}),
	)
	if err != nil {
		panic(err)
	}
	interactClient = c
}

func LikeAction(ctx context.Context, request *interact.LikeRequest) (r *interact.LikeResponse, err error) {
	resp, err := interactClient.LikeAction(ctx, request)
	return resp, err
}
func GetLikeList(ctx context.Context, request *interact.LikeListRequest) (r *interact.LikeListResponse, err error) {
	resp, err := interactClient.GetLikeList(ctx, request)
	return resp, err
}
func CommentAction(ctx context.Context, request *interact.CommentRequest) (r *interact.CommentResponse, err error) {
	resp, err := interactClient.CommentAction(ctx, request)
	return resp, err
}
func GetCommentList(ctx context.Context, request *interact.CommentListRequest) (r *interact.CommentListResponse, err error) {
	resp, err := interactClient.GetCommentList(ctx, request)
	return resp, err
}
