package wsApp

import (
	"fmt"
	"ws-gateway/response"
	"ws-gateway/utils"
)

var (
	closeMessageType       = "close"
	bindGroupMessageType   = "bind_group"
	unBindGroupMessageType = "unbind_group"
)

type wsRequest struct {
	MessageType string      `json:"message_type"`
	Data        interface{} `json:"data"`
}

type groupRequest struct {
	GroupId string `json:"group_id"`
}

func (r *wsRequest) toStruct(msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			var outErr error
			switch x := (err).(type) {
			case error:
				outErr = x
			case string:
				outErr = fmt.Errorf(x)
			case fmt.Stringer:
				outErr = fmt.Errorf(x.String())
			default:
				outErr = fmt.Errorf("%v", x)
			}
			r.MessageType = "error"
			r.Data = utils.StructToBytes(response.WsResponse{
				MessageType: "error",
				Data: response.StandardResponse{
					Code:    100,
					Message: outErr.Error()},
			})
		}
	}()
	err := utils.ByteToStruct(r, msg)
	if err != nil {
		panic(err)
	}
	switch r.MessageType {
	case closeMessageType:
		r.MessageType = closeMessageType
	case bindGroupMessageType:
		var data *wsRequest = &wsRequest{Data: groupRequest{}}
		err := utils.ByteToStruct(data, msg)
		if err != nil {
			panic(err)
		}
		r.MessageType = data.MessageType
		r.Data = data.Data
	case unBindGroupMessageType:
		var data *wsRequest = &wsRequest{Data: groupRequest{}}
		err := utils.ByteToStruct(data, msg)
		if err != nil {
			panic(err)
		}
		r.MessageType = data.MessageType
		r.Data = data.Data
	}
}

type indexRequest struct {
	GroupId string `json:"group_id" form:"group_id"`
	UserId  string `json:"user_id" form:"user_id"`
	Index   int    `json:"index" form:"index" binding:"required"`
}
