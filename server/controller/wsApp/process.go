package wsApp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"ws-gateway/exception"
	"ws-gateway/logger"
	"ws-gateway/redisUtils"
	"ws-gateway/request"
	"ws-gateway/response"
	"ws-gateway/utils"
	"ws-gateway/wsManager"
)

func wsApp(c *gin.Context) {
	var res request.ConnectQuery
	if c.ShouldBindQuery(&res) != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("wsApp.wsApp").
			SetErrorCode(1))
	}
	res.ClientIp = c.ClientIP()
	res.Platform = "app"
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger("wsApp Upgrade", "error", err, "ws升级错误")
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("连接异常").
			SetFunctionName("wsApp.wsApp").
			SetOriginalError(err).
			SetErrorCode(1))
	}
	conn.SetReadLimit(10240)
	initWs(&res, conn)
}

func initWs(r *request.ConnectQuery, conn *websocket.Conn) {
	wsClient := &wsManager.Client{
		ClientId:        wsManager.GenClientId(),
		Socket:          conn,
		ConnectTime:     time.Now(),
		Query:           *r,
		LastMessageTime: time.Now().Unix(),
	}
	wsClient.Read(func(messageType int, message []byte, c *wsManager.Client) error {
		if fmt.Sprintf("%s", message) == "ping" {
			return nil
		}
		var mData wsRequest
		mData.toStruct(message)
		switch mData.MessageType {
		case closeMessageType:
			wsManager.Manager.DisConnect <- wsClient
		case bindGroupMessageType:
			// 绑定群组
			func(res wsRequest, client *wsManager.Client) {
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
						logger.Logger("wsBoardApp Read", "error", outErr, "ws读取错误")
					}
				}()
				groupId := utils.InterfaceToString(res.Data.(map[string]interface{})["group_id"])
				wsManager.Manager.AddClientToGroup(client, groupId)
			}(mData, wsClient)
		case unBindGroupMessageType:
			// 移除组
			func(res wsRequest, client *wsManager.Client) {
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
						logger.Logger("wsBoardApp Read", "error", outErr, "ws读取错误")
					}
				}()
				groupId := utils.InterfaceToString(res.Data.(map[string]interface{})["group_id"])
				wsManager.Manager.DelClientFromGroup(client, groupId)
			}(mData, wsClient)
		}
		return nil
	})
	wsManager.Manager.Connect <- wsClient
}

// groupIndexMsg 组指定序号的消息
// @Summary  组消息
// @Tags 组消息-外部
// @Param Body query indexRequest true "推送消息"
// @Router /api/public/ws-gateway-service/group/msg/index [get]
// @Produce json
// @Success 200 object response.StandardResponse{data=response.GroupPublishData}
func groupIndexMsg(c *gin.Context) {
	var request indexRequest
	if queryErr := c.ShouldBindQuery(&request); queryErr != nil {
		if jsonErr := c.ShouldBindJSON(&request); jsonErr != nil {
			panic(exception.StandardRuntimeBadError().SetOutPutMessage("参数错误").SetFunctionName("wsBoardApp.groupIndexMsg").SetErrorCode(1))
		}
	}
	key := fmt.Sprintf("ws:group:msg:%s:%d", request.GroupId, request.Index)
	data, err := redisUtils.Get(key)
	if err != nil {
		panic(exception.StandardRuntimeBadError().SetOutPutMessage("消息不存在").SetFunctionName("wsBoardApp.groupIndexMsg").SetErrorCode(1))
	}
	var res response.GroupPublishData
	if data != nil {
		_ = json.Unmarshal(data, &res)
	}
	c.JSON(http.StatusOK, &response.StandardResponse{
		Code:    0,
		Message: "success",
		Data:    res,
	})
}

// userIndexMsg 用户指定序号的消息
// @Summary  组消息
// @Tags 组消息-外部
// @Param Body query indexRequest true "推送消息"
// @Router /api/public/ws-gateway-service/user/msg/index [get]
// @Produce json
// @Success 200 object response.StandardResponse{data=response.UserPublishData}
func userIndexMsg(c *gin.Context) {
	var request indexRequest
	if queryErr := c.ShouldBindQuery(&request); queryErr != nil {
		if jsonErr := c.ShouldBindJSON(&request); jsonErr != nil {
			panic(exception.StandardRuntimeBadError().SetOutPutMessage("参数错误").SetFunctionName("wsBoardApp.userIndexMsg").SetErrorCode(1))
		}
	}
	key := fmt.Sprintf("ws:user:msg:%s:%d", request.UserId, request.Index)
	data, err := redisUtils.Get(key)
	if err != nil {
		panic(exception.StandardRuntimeBadError().SetOutPutMessage("消息不存在").SetFunctionName("wsBoardApp.userIndexMsg").SetErrorCode(1))
	}
	var res response.UserPublishData
	if data != nil {
		_ = json.Unmarshal(data, &res)
	}
	c.JSON(http.StatusOK, &response.StandardResponse{
		Code:    0,
		Message: "success",
		Data:    res,
	})
}
