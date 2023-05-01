namespace go base

struct UserInfo{
    1: i64 id
    2: string name
    3: i64 follow_count
    4: i64 follower_count
    5: bool is_follow
}

struct RegisterRequest {
    1: string username
    2: string password
}

struct RegisterResponse {
    1: i64 status_code
    2: string status_msg
    3: i64 user_id
    4: string token
}

struct LoginRequest {
    1: string username (api.query="username")
    2: string password (api.query="password")
}

struct LoginResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional i64 user_id
    4: optional string token
}

struct UserInfoRequest {
    1: i64 to_user_id
    2: string token
    3: optional i64 user_id
    4: optional string username
}

struct UserInfoResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional UserInfo user
}

struct VideoInfo {
    1: i64 id// 视频唯一标识
    2: UserInfo author// 视频作者信息
    3: string play_url// 视频播放地址
    4: string cover_url// 视频封面地址
    5: i64 favorite_count// 视频的点赞总数
    6: i64 comment_count// 视频的评论总数
    7: bool is_favorite// true-已点赞，false-未点赞
    8: string title// 视频标题
}

struct FeedRequest{
    1: optional string latest_time// 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
    2: optional string token
    3: optional i64 user_id
    4: optional string user_name
}

struct FeedResponse{
    1: i64 status_code// 状态码，0-成功，其他值-失败
    2: optional string status_msg// 返回状态描述
    3: optional i64 next_time// 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
    4: optional list<VideoInfo> video_list// 视频列表
}

struct PublishRequest {
    1: binary data
    2: string token
    3: string title
    4: i64 user_id
    5: string user_name
}

struct PublishResponse {
    1: i64 status_code
    2: optional string status_msg
}

struct PublishListRequest {
    1: i64 user_id
    2: i64 query_id
}

struct PublishListResponse {
    1: i64 status_code
    2: optional string status_msg
    4: optional list<VideoInfo> video_list
}

service BaseService {
    RegisterResponse UserRegister(1: RegisterRequest request)(api.post="/douyin/user/register/")
    LoginResponse UserLogin(1: LoginRequest request)(api.post="/douyin/user/login/")
    UserInfoResponse GetUserInfo(1: UserInfoRequest request)(api.get="/douyin/user/")
    FeedResponse GetVideoList(1: FeedRequest request)(api.get="/douyin/feed/")
    PublishResponse PublishAction(1: PublishRequest request) (api.post="/douyin/publish/action/")
    PublishListResponse GetPublishList(1: PublishListRequest request) (api.get="/douyin/publish/list/")
}
