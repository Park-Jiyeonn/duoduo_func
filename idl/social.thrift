namespace go social

include "base.thrift"

struct FollowRequest {
    1: string token
    2: string to_user_id
    3: string action_type
    4: optional i64 user_id
    5: optional string user_name
}

struct Message {
    1: required i64 id;
    2: required i64 to_user_id;
    3: required i64 from_user_id;
    4: required string content;
    5: optional i64 create_time;
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
    3: optional list<base.UserInfo> user_list
}

struct FollowerListRequest {
    1: i64 user_id
    2: string token
    3: optional string user_name
}

struct FollowerListResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional list<base.UserInfo> user_list
}

struct FriendListRequest {
    1: i64 user_id
    2: string token
}

struct FriendListResponse {
    1: i64 status_code
    2: optional string status_msg
    3: optional list<base.UserInfo> user_list
}

struct MessageActionReq {
    1: required string token;
    2: required i64 to_user_id;
    3: required i32 action_type;
    4: required string content;
    5: optional i64 user_id
    6: optional string user_name
}

struct MessageActionResp {
    1: required i32 status_code;
    2: optional string status_msg;
}

struct MessageChatReq {
    1: required string token;
    2: required i64 to_user_id;
    3: i64 pre_msg_time;
    4: optional i64 user_id
    5: optional string user_name
}

struct MessageChatResp {
    1: required i32 status_code;
    2: optional string status_msg;
    3: required list<Message> message_list;
}

service SocialService {
    FollowResponse FollowAction(1: FollowRequest request)(api.post="/douyin/relation/action/")
    FollowingListResponse GetFollowList(1: FollowingListRequest request)(api.get="/douyin/relation/follow/list/")
    FollowerListResponse GetFollowerList(1: FollowerListRequest request)(api.get="/douyin/relation/follower/list/")
    FriendListResponse GetFriendList(1: FriendListRequest request) (api.get="/douyin/relation/friend/list/")
    MessageChatResp MessageChat(1: MessageChatReq req) (api.get="/douyin/message/chat/");
    MessageActionResp MessageAction(1: MessageActionReq req) (api.post="/douyin/message/action/")
}