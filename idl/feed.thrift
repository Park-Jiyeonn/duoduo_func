namespace go feed

include "user.thrift"

struct VideoInfo {
    1: i64 id// 视频唯一标识
    2: user.UserInfo author// 视频作者信息
    3: string play_url// 视频播放地址
    4: string cover_url// 视频封面地址
    5: i64 favorite_count// 视频的点赞总数
    6: i64 comment_count// 视频的评论总数
    7: bool is_favorite// true-已点赞，false-未点赞
    8: string title// 视频标题
}

struct FeedRequest{
    1: optional string latest_time// 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
    2: optional i64 user_id
}

struct FeedResponse{
    1: i64 status_code// 状态码，0-成功，其他值-失败
    2: optional string status_msg// 返回状态描述
    3: optional i64 next_time// 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
    4: optional list<VideoInfo> video_list// 视频列表
}

struct GetByIDRequest{
    1: required i64 query_id
    2: required i64 user_id
}

service FeedService{
    FeedResponse GetVideoList(1: FeedRequest request)(api.get="/douyin/feed/")
    FeedResponse GetVideoListById(1: GetByIDRequest request)
}
