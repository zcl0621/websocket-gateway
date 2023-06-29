package wsManager

import (
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"sync/atomic"
	"time"
	"ws-gateway/logger"
	"ws-gateway/request"
)

type Client struct {
	ClientId        string          // 标识ID
	Socket          *websocket.Conn // 用户连接
	ConnectTime     time.Time       // 首次连接时间
	IsDeleted       bool            // 是否删除或下线
	Query           request.ConnectQuery
	LastMessageTime int64
}

func GenClientId() string {
	return xid.New().String()
}

func (c *Client) Read(handlerFun func(messageType int, message []byte, c *Client) error) {
	go func() {
		for {
			messageType, message, err := c.Socket.ReadMessage()
			if err != nil {
				if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					Manager.DisConnect <- c
					return
				} else if messageType != websocket.PingMessage {
					return
				}
			} else if messageType == websocket.TextMessage {
				now := time.Now().Unix()
				atomic.SwapInt64(&c.LastMessageTime, now)
				if string(message) == "ping" {
					ToClientChan <- &clientInfo{
						Client:      c,
						Msg:         []byte("pong"),
						MessageType: websocket.TextMessage,
					}
				} else if string(message) == "Ping" {
					ToClientChan <- &clientInfo{
						Client:      c,
						Msg:         []byte("Pong"),
						MessageType: websocket.TextMessage,
					}
				} else if string(message) == "pong" {
					ToClientChan <- &clientInfo{
						Client:      c,
						Msg:         []byte("ping"),
						MessageType: websocket.TextMessage,
					}
				} else if string(message) == "Pong" {
					ToClientChan <- &clientInfo{
						Client:      c,
						Msg:         []byte("Ping"),
						MessageType: websocket.TextMessage,
					}
				} else {
					sErr := handlerFun(messageType, message, c)
					if sErr != nil {
						logger.Logger("client Read", "waring", sErr, "")
						return
					}
				}
			} else if messageType == websocket.PingMessage {
				ToClientChan <- &clientInfo{
					Client:      c,
					Msg:         []byte{},
					MessageType: websocket.PongMessage,
				}
			} else if messageType == websocket.PongMessage {
				ToClientChan <- &clientInfo{
					Client:      c,
					Msg:         []byte{},
					MessageType: websocket.PingMessage,
				}
			}
		}
	}()

}
