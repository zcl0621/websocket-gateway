package wsManager

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"ws-gateway/logger"
	"ws-gateway/response"
	"ws-gateway/shardMap"
)

// 连接管理
type ClientManager struct {
	ClientIdMap     shardMap.ConcurrentMap[string, *Client]         // 保存连接id和链接的映射
	UserIdMap       shardMap.ConcurrentMap[string, map[string]bool] // 保存用户id和链接的映射 一个用户可以有多个链接 key 用户id value map[链接id]bool
	userIdLock      *sync.RWMutex
	GroupIdMap      shardMap.ConcurrentMap[string, map[string]bool] // 保存组id和链接的映射 一个组可以有多个链接 key 组id value map[链接id]bool
	groupIdLock     *sync.RWMutex
	ClientGroupMap  shardMap.ConcurrentMap[string, map[string]bool] // 保存连接id和组id的映射 一个连接可以有多个组 key 连接id value map[组id]bool
	clientGroupLock *sync.RWMutex

	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开连接处理
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		ClientIdMap:     shardMap.New[*Client](),
		UserIdMap:       shardMap.New[map[string]bool](),
		userIdLock:      &sync.RWMutex{},
		GroupIdMap:      shardMap.New[map[string]bool](),
		groupIdLock:     &sync.RWMutex{},
		ClientGroupMap:  shardMap.New[map[string]bool](),
		clientGroupLock: &sync.RWMutex{},
		Connect:         make(chan *Client, 2048),
		DisConnect:      make(chan *Client, 2048),
	}

	return
}

// 管道处理程序
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			// 建立连接事件
			manager.EventConnect(client)
		case conn := <-manager.DisConnect:
			// 断开连接事件
			manager.EventDisconnect(conn)
		}
	}
}

// 建立连接事件
func (manager *ClientManager) EventConnect(client *Client) {
	manager.AddClient(client)
	GroupChannelsManager.ClientStatusProducts <- &WsClientStatusData{
		ClientId: client.ClientId,
		UserId:   client.Query.UserId,
		Type:     "enter",
		Platform: client.Query.Platform,
		OS:       client.Query.OS,
		Version:  client.Query.Version,
		IP:       client.Query.ClientIp,
		AppType:  client.Query.AppType,
	}
}

// 断开连接时间
func (manager *ClientManager) EventDisconnect(client *Client) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("client EventDisconnect Error", logger.ERROR, nil, fmt.Sprintf("err: %v", err))
		}
	}()
	logger.Logger("client EventDisconnect", logger.INFO, nil, fmt.Sprintf("连接已断开 clientId:%s", client.ClientId))
	//关闭连接
	err := client.Socket.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "close"), time.Now().Add(time.Second))
	if err != nil && err != websocket.ErrCloseSent {
		// If close message could not be sent, then close without the handshake.
		_ = client.Socket.Close()
	}
	manager.DelClient(client)
	client.IsDeleted = true

	GroupChannelsManager.ClientStatusProducts <- &WsClientStatusData{
		ClientId: client.ClientId,
		UserId:   client.Query.UserId,
		Type:     "leave",
		Platform: client.Query.Platform,
	}

	client = nil
}

// AddClientToClientIdMap 添加客户端到ClientIdMap
func (manager *ClientManager) AddClientToClientIdMap(client *Client) {
	manager.ClientIdMap.Set(client.ClientId, client)
}

// GetClientByClientId 通过ClientId获取客户端
func (manager *ClientManager) GetClientByClientId(clientId string) (*Client, bool) {
	if client, ok := manager.ClientIdMap.Get(clientId); !ok {
		return nil, false
	} else {
		return client, true
	}
}

// DelClientFromClientIdMap 从ClientIdMap删除客户端
func (manager *ClientManager) DelClientFromClientIdMap(client *Client) {
	manager.ClientIdMap.Remove(client.ClientId)
}

// AddClientToUserIdMap 添加客户端到UserIdMap
func (manager *ClientManager) AddClientToUserIdMap(client *Client) {
	if client.Query.UserId != "" {
		manager.userIdLock.Lock()
		defer manager.userIdLock.Unlock()
		if userIdMap, ok := manager.UserIdMap.Get(client.Query.UserId); ok {
			userIdMap[client.ClientId] = true
			manager.UserIdMap.Set(client.Query.UserId, userIdMap)
		} else {
			manager.UserIdMap.Set(client.Query.UserId, map[string]bool{client.ClientId: true})
		}
	}
}

// DelClientFromUserIdMap 从UserIdMap删除客户端
func (manager *ClientManager) DelClientFromUserIdMap(client *Client) {
	if client.Query.UserId != "" {
		manager.userIdLock.Lock()
		defer manager.userIdLock.Unlock()
		if userIdMap, ok := manager.UserIdMap.Get(client.Query.UserId); ok {
			delete(userIdMap, client.ClientId)
			manager.UserIdMap.Set(client.Query.UserId, userIdMap)
		}
	}
}

func (manager *ClientManager) AddClientToGroup(client *Client, groupId string) {
	manager.AddClientToGroupIdMap(client, groupId)
	manager.AddClientToClientGroupMap(client, groupId)
	GroupChannelsManager.ClientStatusProducts <- &WsClientStatusData{
		GroupId:  groupId,
		ClientId: client.ClientId,
		UserId:   client.Query.UserId,
		Type:     "group_enter",
		Platform: client.Query.Platform,
	}
}

// AddClientToGroupIdMap 添加客户端到GroupIdMap
func (manager *ClientManager) AddClientToGroupIdMap(client *Client, groupId string) {
	if groupId != "" {
		manager.groupIdLock.Lock()
		defer manager.groupIdLock.Unlock()
		if groupIdMap, ok := manager.GroupIdMap.Get(groupId); ok {
			groupIdMap[client.ClientId] = true
			manager.GroupIdMap.Set(groupId, groupIdMap)
		} else {
			manager.GroupIdMap.Set(groupId, map[string]bool{client.ClientId: true})
		}
	}
}

// DelClientFromGroupIdMap 从GroupIdMap删除客户端
func (manager *ClientManager) DelClientFromGroupIdMap(client *Client, groupId string) {
	if groupId != "" {
		manager.groupIdLock.Lock()
		defer manager.groupIdLock.Unlock()
		if groupIdMap, ok := manager.GroupIdMap.Get(groupId); ok {
			delete(groupIdMap, client.ClientId)
			manager.GroupIdMap.Set(groupId, groupIdMap)
			GroupChannelsManager.ClientStatusProducts <- &WsClientStatusData{
				GroupId:  groupId,
				ClientId: client.ClientId,
				UserId:   client.Query.UserId,
				Type:     "group_leave",
				Platform: client.Query.Platform,
			}
		}
	}
}

// AddClientToClientGroupMap 添加客户端到ClientGroupMap
func (manager *ClientManager) AddClientToClientGroupMap(client *Client, groupId string) {
	if groupId != "" {
		manager.clientGroupLock.Lock()
		defer manager.clientGroupLock.Unlock()
		if groupClients, ok := manager.ClientGroupMap.Get(client.ClientId); ok {
			groupClients[groupId] = true
			manager.ClientGroupMap.Set(client.ClientId, groupClients)
		} else {
			manager.ClientGroupMap.Set(client.ClientId, map[string]bool{groupId: true})
		}
		logger.Logger("addClientToGroup", "info", nil, fmt.Sprintf("client_id: %s 新增组 group_id: %s", client.ClientId, groupId))
	}
}

func (manager *ClientManager) DelClientFromGroup(client *Client, groupId string) {
	manager.DelClientFromClientGroupMap(client, groupId)
	manager.DelClientFromGroupIdMap(client, groupId)

}

// DelClientFromClientGroupMap 从ClientGroupMap删除客户端
func (manager *ClientManager) DelClientFromClientGroupMap(client *Client, groupId string) {
	if groupId != "" {
		manager.clientGroupLock.Lock()
		defer manager.clientGroupLock.Unlock()
		if groupClients, ok := manager.ClientGroupMap.Get(client.ClientId); ok {
			delete(groupClients, groupId)
			manager.ClientGroupMap.Set(client.ClientId, groupClients)
		}
		logger.Logger("delClientFromGroup", "info", nil, fmt.Sprintf("client_id: %s 删除组 group_id: %s", client.ClientId, groupId))
	}
}

// AddClient 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	logger.Logger("AddClient", "info", nil, fmt.Sprintf("client_id: %s 新增连接", client.ClientId))
	manager.AddClientToClientIdMap(client)
	manager.AddClientToUserIdMap(client)
	clientInfo := response.WsResponse{
		MessageType: "client_info",
		Data:        map[string]interface{}{"client_id": client.ClientId},
	}
	d, _ := json.Marshal(clientInfo)
	SendMessage(client, d, websocket.TextMessage)
}

// DelClient 删除客户端
func (manager *ClientManager) DelClient(client *Client) {
	logger.Logger("DelClient", "info", nil, fmt.Sprintf("client_id: %s 断开连接", client.ClientId))
	manager.DelClientFromClientIdMap(client)
	manager.DelClientFromUserIdMap(client)
	// 删除客户端的所有组
	if groupClients, ok := manager.ClientGroupMap.Get(client.ClientId); ok {
		for groupId := range groupClients {
			manager.DelClientFromGroupIdMap(client, groupId)
		}
	}
	manager.ClientGroupMap.Remove(client.ClientId)
}

// AllClient 获取所有客户端
func (manager *ClientManager) AllClient() map[string]*Client {
	return manager.ClientIdMap.Items()
}

// SendToClient 发送消息给客户端
func (manager *ClientManager) SendToClient(clientId string, msg []byte) {
	if client, ok := manager.ClientIdMap.Get(clientId); ok {
		SendMessage(client, msg, websocket.TextMessage)
	}
}

// SendToUser 发送消息给用户
func (manager *ClientManager) SendToUser(userId string, msg []byte) {
	if userIdMap, ok := manager.UserIdMap.Get(userId); ok {
		for clientId := range userIdMap {
			manager.SendToClient(clientId, msg)
		}
	}
}

// SendToGroup 发送消息给组
func (manager *ClientManager) SendToGroup(groupId string, msg []byte) {
	if groupIdMap, ok := manager.GroupIdMap.Get(groupId); ok {
		for clientId := range groupIdMap {
			manager.SendToClient(clientId, msg)
		}
	}
}
