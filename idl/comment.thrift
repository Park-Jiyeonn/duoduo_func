namespace go comment

include "user.thrift"

struct CommentInfo {
    1: i64 id
    2: user.UserInfo user
    3: string content
    4: string create_date
}

struct CommentRequest {
    1: string token
    2: i64 video_id
    3: bool action_type
    4: optional string comment_text
    5: optional i64 comment_id
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

service CommentService {
    CommentResponse CommentAction(1: CommentRequest request)(api.post="/douyin/comment/action/")
    CommentListResponse GetCommentList(1: CommentListRequest request)(api.get="/douyin/comment/list/")
}