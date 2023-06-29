package wsManager

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
	"ws-gateway/logger"
)

var ToClientChan chan *clientInfo

// channel通道结构体
type clientInfo struct {
	Client      *Client
	Msg         []byte
	MessageType int
}

// 发送信息
func SendMessage(clientId *Client, msg []byte, messageType int) {
	ToClientChan <- &clientInfo{Client: clientId, Msg: msg, MessageType: messageType}
	return
}

// 监听并发送给客户端信息
func WriteMessage() {
	for {
		ci := <-ToClientChan
		func(ci *clientInfo) {
			Render(ci.Client, ci.Msg, ci.MessageType)
		}(ci)
	}
}

func Render(c *Client, msg []byte, messageType int) {
	var err error
	defer func() {
		if rErr := recover(); rErr != nil {
			err = errors.New(fmt.Sprintf("%s", rErr))
			logger.Logger("websocket Render Defer func", logger.ERROR, err, fmt.Sprintf("e: %s", err.Error()))
		}
	}()
	err = c.Socket.SetWriteDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		logger.Logger("websocket Render", logger.ERROR, err, fmt.Sprintf("clientId: %s msg: %s e: %s", c.ClientId, msg, err.Error()))
		Manager.DisConnect <- c
		return
	}
	switch messageType {
	case websocket.TextMessage:
		err = c.Socket.WriteMessage(websocket.TextMessage, msg)
		break
	case websocket.PingMessage:
		err = c.Socket.WriteMessage(websocket.PingMessage, []byte{})
		break
	case websocket.PongMessage:
		err = c.Socket.WriteMessage(websocket.PongMessage, []byte{})
		break
	default:
		err = c.Socket.WriteMessage(websocket.TextMessage, msg)
		break
	}
	if err != nil {
		logger.Logger("Render 发消息失败 正在断开连接", logger.ERROR, nil, fmt.Sprintf("clientId: %s e: %s", c.ClientId, err.Error()))
		Manager.DisConnect <- c
	}
	return
}
