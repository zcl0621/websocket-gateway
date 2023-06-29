package response

type StandardResponse struct {
	Code    int         `json:"code"`    //0 成功 1失败
	Message string      `json:"message"` //提示语
	Data    interface{} `json:"data"`
}

type GroupPublishData struct {
	GroupId string `json:"group_id"` //消息组
	Message string `json:"message"`  //消息数据 json字符串
	Index   int    `json:"index"`    //消息索引 提交消息时自增
}

type UserPublishData struct {
	UserId  string `json:"user_id"` //用户id
	Message string `json:"message"` //消息数据 json字符串
	Index   int    `json:"index"`   //消息索引 提交消息时自增
}

type ClientPublishData struct {
	ClientId string `json:"client_id"` //客户端id
	Message  string `json:"message"`   //消息数据 json字符串
	Index    int    `json:"index"`     //消息索引 提交消息时自增
}

type WsResponse struct {
	MessageType string      `json:"message_type"`
	Data        interface{} `json:"data"`
}
