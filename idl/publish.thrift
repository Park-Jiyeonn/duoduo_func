namespace go publish

include "feed.thrift"

struct PublishRequest {
    1: binary data
    2: i64 user_id
    3: string title
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
    4: optional list<feed.VideoInfo> video_list
}

service PublishService {
    PublishResponse PublishAction(1: PublishRequest request) (api.post="/douyin/publish/action/")
    PublishListResponse GetPublishList(1: PublishListRequest request) (api.get="/douyin/publish/list/")
}