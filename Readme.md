# WebSocket接口

# API: /api/public/ws-gateway-service/ws

> 支持Token和AnonymousToken两种方式

## 信息返回

### 首次返回

```json
{
  "message_type": "client_info",
  "data": {
    "client_id": "xxxx"
  }
}
```

## 绑定消息组

```
{
    "message_type":"bind_group",
    "group_id": "test_1"
}
```

## 取消绑定消息组

```
{
    "message_type":"unbind_group",
    "group_id": "test_1"
}
```

# 内部接口

## 组发消息

### url: /api/inner/ws-holder-service/group-publish

### method: POST

### body:

```json
{
  "group_id": "test_1",
  //组id
  "message": "{\"name\":\"test\",\"group_id\":\"test_2\"}"
  //消息内容 json格式 须包含group_id
}
```

## 用户消息

### url: /api/inner/ws-holder-service/user-publish

### method: POST

### body:

```json
{
  "user_id": "test_1",
  //用户id
  "message": "{\"name\":\"test\",\"group_id\":\"test_2\"}"
  //消息内容 json格式 须包含group_id
}
```

## 链接消息

### url: /api/inner/ws-holder-service/client-publish

### method: POST

### body

```json
{
  "client_id": "test_1",
  //链接id
  "message": "{\"name\":\"test\",\"group_id\":\"test_2\"}"
  //消息内容 json格式 须包含group_id
}
```

# kafka

## 链接事件

### topic：ty-ws-client-status

### message
```go

package main

type WSClientStatusKafkaMessage struct {
    ClientId string `json:"client_id"` // 客户端id
    UserId   string `json:"user_id"`   //用户id
    UserType string `json:"user_type"` // user guest
    Type     string `json:"type"`      // enter leave
}
```

## 组链接变化事件

### topic: ty-ws-group-status

### message:

```go

package main

type WSGroupStatusKafkaMessage struct {
	GroupId  string `json:"group_id"`  // 组id
	ClientId string `json:"client_id"` // 客户端id
	UserId   string `json:"user_id"`   //用户id
	UserType string `json:"user_type"` // user guest
	Type     string `json:"type"`      // enter leave
}

```


# 约定

## 组

### 待补充