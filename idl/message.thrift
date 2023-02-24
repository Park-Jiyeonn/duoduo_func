struct Message {
    1: required i64 id;
    2: required i64 to_user_id;
    3: required i64 from_user_id;
    4: required string content;
    5: optional i64 create_time;
}

struct MessageActionReq {
    1: required string token;
    2: required i64 to_user_id;
    3: required i32 action_type;
    4: required string content;
}

struct MessageActionResp {
    1: required i32 status_code;
    2: optional string status_msg;
}

struct MessageChatReq {
    1: required string token;
    2: required i64 to_user_id;
    3: i64 pre_msg_time;
}

struct MessageChatResp {
    1: required i32 status_code;
    2: optional string status_msg;
    3: required list<Message> message_list;
}

service MessageService {
    MessageChatResp MessageChat(1: MessageChatReq req) (api.get="/douyin/message/chat/");
    MessageActionResp MessageAction(1: MessageActionReq req) (api.post="/douyin/message/action/")
}