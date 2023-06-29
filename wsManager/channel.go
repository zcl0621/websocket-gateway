package wsManager

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"ws-gateway/logger"
	"ws-gateway/redisUtils"
	"ws-gateway/response"
)

var redisGroupChannels = "ws-gateway-channels:group"
var redisUserChannels = "ws-gateway-channels:user"
var redisClientChannels = "ws-gateway-channels:client"
var redisClientStatusChannels = "ws-gateway-channels:client-status"
var redisClientStatusSteamChannels = "ws-gateway-channels:client-steam-status"

type WsClientStatusData struct {
	GroupId  string `json:"group_id"`
	ClientId string `json:"client_id"`
	UserId   string `json:"user_id"`
	Type     string `json:"type"`     // enter leave group_enter group_leave
	Platform string `json:"platform"` // board app
	OS       string `json:"os"`
	Version  string `json:"version"`
	IP       string `json:"ip"`
	AppType  string `json:"app_type"` // higo airschool
}

type Channels struct {
	GroupProducts        chan *response.GroupPublishData
	GroupCustomer        chan *response.GroupPublishData
	UserProducts         chan *response.UserPublishData
	UserCustomer         chan *response.UserPublishData
	ClientCustomer       chan *response.ClientPublishData
	ClientProducts       chan *response.ClientPublishData
	ClientStatusProducts chan *WsClientStatusData
}

func NewGroupChannelsManager() (channelsManager *Channels) {
	channelsManager = &Channels{
		GroupProducts:        make(chan *response.GroupPublishData, 2048),
		GroupCustomer:        make(chan *response.GroupPublishData, 2048),
		UserProducts:         make(chan *response.UserPublishData, 2048),
		UserCustomer:         make(chan *response.UserPublishData, 2048),
		ClientProducts:       make(chan *response.ClientPublishData, 2048),
		ClientCustomer:       make(chan *response.ClientPublishData, 2048),
		ClientStatusProducts: make(chan *WsClientStatusData, 2048),
	}
	return channelsManager
}

// 管道处理程序
func (c *Channels) Start() {
	for {
		select {
		case data := <-c.GroupProducts:
			//组消息接收
			c.handlerGroupProducts(data)
		case data := <-c.GroupCustomer:
			//组消息发送
			c.handlerGroupCustomer(data)
		case data := <-c.UserProducts:
			//用户消息接收
			c.handlerUserProducts(data)
		case data := <-c.UserCustomer:
			//用户消息发送
			c.handlerUserCustomer(data)
		case data := <-c.ClientProducts:
			c.handlerClientProducts(data)
		case data := <-c.ClientCustomer:
			c.handlerClientCustomer(data)
		case data := <-c.ClientStatusProducts:
			c.handlerClientStatusProducts(data)
		}
	}
}

func (c *Channels) handlerClientStatusProducts(data *WsClientStatusData) {
	bData, mErr := json.Marshal(data)
	if mErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't marshal data value:%v", data))
		return
	}
	rErr := redisUtils.Publish(redisClientStatusChannels, fmt.Sprintf("%s", bData))
	if rErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't Publish data value:%s", bData))
	}
	sErr := redisUtils.XAdd(redisClientStatusSteamChannels, bData)
	if sErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't Steam data value:%s", bData))
	}
}

func (c *Channels) handlerGroupProducts(data *response.GroupPublishData) {
	bData, mErr := json.Marshal(data)
	if mErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't marshal data value:%v", data))
		return
	}
	rErr := redisUtils.Publish(redisGroupChannels, fmt.Sprintf("%s", bData))
	if rErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't Publish data value:%s", bData))

	}
}

func (c *Channels) handlerUserProducts(data *response.UserPublishData) {
	bData, mErr := json.Marshal(data)
	if mErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't marshal data value:%v", data))
		return
	}
	rErr := redisUtils.Publish(redisUserChannels, fmt.Sprintf("%s", bData))
	if rErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't Publish data value:%s", bData))

	}
}

func (c *Channels) handlerClientProducts(data *response.ClientPublishData) {
	bData, mErr := json.Marshal(data)
	if mErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't marshal data value:%v", data))
		return
	}
	rErr := redisUtils.Publish(redisClientChannels, fmt.Sprintf("%s", bData))
	if rErr != nil {
		logger.Logger("channel.handlerGroupProducts", logger.ERROR, mErr, fmt.Sprintf("can't Publish data value:%s", bData))

	}
}

func (c *Channels) handlerClientCustomer(data *response.ClientPublishData) {
	bMessage, _ := json.Marshal(data)
	Manager.SendToClient(data.ClientId, bMessage)
}

func (c *Channels) handlerGroupCustomer(data *response.GroupPublishData) {
	bMessage, _ := json.Marshal(data)
	Manager.SendToGroup(data.GroupId, bMessage)
}

func (c *Channels) handlerUserCustomer(data *response.UserPublishData) {
	bMessage, _ := json.Marshal(data)
	Manager.SendToUser(data.UserId, bMessage)
}

func SubRedis() {
	go func() {
		_, err := redisUtils.Subscribe(redisGroupChannels, func(message *redis.Message, err error) {
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't subscribe message value:%v", message))
				return
			}
			var data response.GroupPublishData
			err = json.Unmarshal([]byte(message.Payload), &data)
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't Unmarshal message value:%v", message))
				return
			}
			GroupChannelsManager.GroupCustomer <- &data
		})
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		_, err := redisUtils.Subscribe(redisUserChannels, func(message *redis.Message, err error) {
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't subscribe message value:%v", message))
				return
			}
			var data response.UserPublishData
			err = json.Unmarshal([]byte(message.Payload), &data)
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't Unmarshal message value:%v", message))
				return
			}
			GroupChannelsManager.UserCustomer <- &data
		})
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		_, err := redisUtils.Subscribe(redisClientChannels, func(message *redis.Message, err error) {
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't subscribe message value:%v", message))
				return
			}
			var data response.ClientPublishData
			err = json.Unmarshal([]byte(message.Payload), &data)
			if err != nil {
				logger.Logger("channel.SubRedis", logger.ERROR, err, fmt.Sprintf("can't Unmarshal message value:%v", message))
				return
			}
			GroupChannelsManager.ClientCustomer <- &data
		})
		if err != nil {
			panic(err)
		}
	}()
}
