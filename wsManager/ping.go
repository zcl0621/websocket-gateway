package wsManager

import (
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

var heartbeatInterval = 8 * time.Second

// 启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for _, conn := range Manager.AllClient() {
				now := time.Now().Unix()
				if now-atomic.LoadInt64(&conn.LastMessageTime) > 2*int64(heartbeatInterval.Seconds()) {
					Manager.DisConnect <- conn
					continue
				}
				ToClientChan <- &clientInfo{
					Client:      conn,
					Msg:         []byte("Ping"),
					MessageType: websocket.TextMessage,
				}
			}
		}
	}()
}
