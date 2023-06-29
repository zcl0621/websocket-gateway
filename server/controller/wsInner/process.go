package wsInner

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ws-gateway/exception"
	"ws-gateway/redisUtils"
	"ws-gateway/response"
	"ws-gateway/wsManager"
)

// groupPublish 推送组消息
// @Summary  推送组消息
// @Tags 推送消息-内部
// @Param data body response.GroupPublishData true "推送消息"
// @Router /api/inner/ws-gateway-service/group-publish [post]
// @Produce json
// @Success 200 object response.StandardResponse{data=publishResponse}
func groupPublish(c *gin.Context) {
	var request response.GroupPublishData
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().SetOutPutMessage("参数错误").SetFunctionName("wsInner.groupPublish").SetErrorCode(1))
	}
	msgIndex, _ := redisUtils.INCRWithExpire(fmt.Sprintf("ws:group:index:%s", request.GroupId), 3600)
	request.Index = msgIndex
	m, _ := json.Marshal(request)
	redisUtils.Set(fmt.Sprintf("ws:group:msg:%s:%d", request.GroupId, msgIndex), m, 120)
	wsManager.GroupChannelsManager.GroupProducts <- &request
	c.JSON(http.StatusOK, &response.StandardResponse{
		Code:    0,
		Message: "success",
		Data: publishResponse{
			Index: msgIndex,
		},
	})
}

// userPublish 推送用户消息
// @Summary  推送用户消息
// @Tags 推送消息-内部
// @Param data body response.UserPublishData true "推送消息"
// @Router /api/inner/ws-gateway-service/user-publish [post]
// @Produce json
// @Success 200 object response.StandardResponse{data=publishResponse}
func userPublish(c *gin.Context) {
	var request response.UserPublishData
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().SetOutPutMessage("参数错误").SetFunctionName("wsInner.userPublish").SetErrorCode(1))
	}
	msgIndex, _ := redisUtils.INCRWithExpire(fmt.Sprintf("ws:user:index:%s", request.UserId), 3600)
	request.Index = msgIndex
	m, _ := json.Marshal(request)
	redisUtils.Set(fmt.Sprintf("ws:user:msg:%s:%d", request.UserId, msgIndex), m, 120)
	wsManager.GroupChannelsManager.UserProducts <- &request
	c.JSON(http.StatusOK, &response.StandardResponse{
		Code:    0,
		Message: "success",
		Data: publishResponse{
			Index: msgIndex,
		},
	})
}

// clientPublish 推送链接消息
// @Summary  推送链接消息
// @Tags 推送消息-内部
// @Param data body response.ClientPublishData true "推送消息"
// @Router /api/inner/ws-gateway-service/client-publish [post]
// @Produce json
// @Success 200 object response.StandardResponse{data=publishResponse}
func clientPublish(c *gin.Context) {
	var request response.ClientPublishData
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().SetOutPutMessage("参数错误").SetFunctionName("wsInner.userPublish").SetErrorCode(1))
	}
	msgIndex, _ := redisUtils.INCRWithExpire(fmt.Sprintf("ws:client:index:%s", request.ClientId), 3600)
	request.Index = msgIndex
	if _, ok := wsManager.Manager.GetClientByClientId(request.ClientId); ok {
		m, _ := json.Marshal(request)
		wsManager.Manager.SendToClient(request.ClientId, m)
	} else {
		wsManager.GroupChannelsManager.ClientProducts <- &request
	}
	c.JSON(http.StatusOK, &response.StandardResponse{
		Code:    0,
		Message: "success",
		Data: publishResponse{
			Index: msgIndex,
		},
	})
}
