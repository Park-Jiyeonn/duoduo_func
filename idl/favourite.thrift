namespace go favorite

include "feed.thrift"

struct LikeRequest {
    1: string token
    2: string video_id
    3: i64 action_type
}

struct LikeResponse {
    1: i64 status_code
    2: string status_msg
}

struct LikeListRequest {
    1: i64 user_id
    2: i64 query_id
}

struct LikeListResponse {
    1: i64 status_code
    2: optional string status_msg
    4: optional list<feed.VideoInfo> video_list
}

service FavoriteService {
    LikeResponse LikeAction(1: LikeRequest request)(api.post="/douyin/favorite/action/")
    LikeListResponse GetLikeList(1: LikeListRequest request)(api.get="/douyin/favorite/list/")
}