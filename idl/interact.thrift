namespace go interact

include "base.thrift"

struct LikeRequest {
    1: string token
    2: i64 video_id
    3: string action_type
    4: optional i64 user_id
    5: optional string user_name
}

struct LikeResponse {
    1: i64 status_code
    2: string status_msg
}

struct LikeListRequest {
    1: i64 user_id
    2: string token
}

struct LikeListResponse {
    1: i64 status_code
    2: optional string status_msg
    4: optional list<base.VideoInfo> video_list
}

struct CommentInfo {
    1: i64 id
    2: base.UserInfo user
    3: string content
    4: string create_date
}

struct CommentRequest {
    1: string token
    2: i64 video_id
    3: string action_type
    4: optional string comment_text
    5: optional i64 comment_id
    6: optional i64 user_id
    7: optional string user_name
}

struct CommentResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional CommentInfo comment
}

struct CommentListRequest {
    1: string token
    2: i64 video_id
}

struct CommentListResponse {
    1: i64 status_code
    2: optional string status_msg
    4: optional list<CommentInfo> comment_list
}

service InteractService {
    LikeResponse LikeAction(1: LikeRequest request)(api.post="/douyin/favorite/action/")
    LikeListResponse GetLikeList(1: LikeListRequest request)(api.get="/douyin/favorite/list/")
    CommentResponse CommentAction(1: CommentRequest request)(api.post="/douyin/comment/action/")
    CommentListResponse GetCommentList(1: CommentListRequest request)(api.get="/douyin/comment/list/")
}