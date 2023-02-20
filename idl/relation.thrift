namespace go relation

include "user.thrift"

struct FollowRequest {
    1: string token
    2: string to_user_id
    3: string action_type
}

struct FollowResponse {
    1: i64 status_code
    2: string status_msg
}

struct FollowingListRequest {
    1: i64 user_id
    2: string token
}

struct FollowingListResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional list<user.UserInfo> user_list
}

struct FollowerListRequest {
    1: i64 user_id
    2: string token
}

struct FollowerListResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional list<user.UserInfo> user_list
}

service RelationService {
    FollowResponse FollowAction(1: FollowRequest request)(api.post="/douyin/relation/action/")
    FollowingListResponse GetFollowList(1: FollowingListRequest request)(api.get="/douyin/relation/follow/list/")
    FollowerListResponse GetFollowerList(1: FollowerListRequest request)(api.get="/douyin/relation/follower/list/")
}