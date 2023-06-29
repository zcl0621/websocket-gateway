package redisUtils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// Publish 发布消息
func Publish(channel string, value string) error {
	return Pool.Publish(context.Background(), channel, value).Err()
}

// Subscribe 订阅
func Subscribe(channel string, handlerMessage func(*redis.Message, error)) (*redis.PubSub, error) {
	pubsub := Pool.Subscribe(context.Background(), channel)
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		handlerMessage(msg, err)
	}
}

// Exists 判断key是否存在
func Exists(key string) (bool, error) {
	res := Pool.Do(context.Background(), "EXISTS", key)
	return res.Bool()
}

// Get 获取key值
func Get(key string) ([]byte, error) {
	res := Pool.Get(context.Background(), key)
	return res.Bytes()
}

// Set 设置key值加过期时间 0 为永久
func Set(key string, value []byte, survivalTime int) error {
	res := Pool.Set(context.Background(), key, value, time.Duration(survivalTime)*time.Second)
	return res.Err()
}

// Del 删除key值
func Del(key string) error {
	res := Pool.Del(context.Background(), key)
	return res.Err()
}

// HSet 向名称为key的hash中添加元素field
func HSet(key string, field string, value []byte, survivalTime int) error {
	res := Pool.HSetNX(context.Background(), key, field, value)
	if res.Err() != nil {
		return res.Err()
	}
	if survivalTime > 0 {
		res = Pool.Expire(context.Background(), key, time.Duration(survivalTime)*time.Second)
		if res.Err() != nil {
			return res.Err()
		}
	}
	return nil
}

// HGETALL 获取名称为key的hash中所有的键（field）及其对应的value
func HGETALL(key string) (map[string]string, error) {
	res := Pool.HGetAll(context.Background(), key)
	return res.Val(), res.Err()
}

// HGETField 获取名称为key的hash中键为field的值
func HGETField(key string, field string) ([]byte, error) {
	return Pool.HGet(context.Background(), key, field).Bytes()
}

// HDelField 删除名称为key的hash中键为field的域
func HDelField(key string, field string) error {
	return Pool.HDel(context.Background(), key, field).Err()
}

// HDel 删除名称为key的hash中所有的数据
func HDel(key string) error {
	return Pool.HDel(context.Background(), key).Err()
}

// LPush 向名称为key的list尾添加一个值为value的 元素
func LPush(key string, value []byte) error {
	return Pool.LPush(context.Background(), key, value).Err()
}

// LRPop 从名称为key的list中取出第一个元素,并将此元素从list中删除
func LRPop(key string) ([]byte, error) {
	return Pool.RPop(context.Background(), key).Bytes()
}

// INCR 递增
func INCR(key string) (int, error) {
	res := Pool.Incr(context.Background(), key)
	return int(res.Val()), res.Err()
}

// INCRWithExpire 递增并设置过期时间
func INCRWithExpire(key string, survivalTime int) (int, error) {
	res := Pool.Incr(context.Background(), key)
	if res.Err() != nil {
		return 0, res.Err()
	}
	tRes := Pool.Expire(context.Background(), key, time.Duration(survivalTime)*time.Second)
	if tRes.Err() != nil {
		return 0, tRes.Err()
	}
	return int(res.Val()), res.Err()
}

// DECR 递减
func DECR(key string) (int, error) {
	res := Pool.Decr(context.Background(), key)
	return int(res.Val()), res.Err()
}

// XAdd 向名称为key的stream中添加一个元素
func XAdd(key string, value []byte) error {
	res := Pool.XAdd(context.Background(), &redis.XAddArgs{
		Stream: key,
		Values: map[string]interface{}{"value": value},
	})
	return res.Err()
}

// XReadGroup 从名称为key的stream中读取数据
func XReadGroup(key string, group string, consumer string) ([]byte, error) {
	res := Pool.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{key, ">"},
		Count:    1,
		Block:    0,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if len(res.Val()) == 0 {
		return nil, nil
	}
	r := []byte(res.Val()[0].Messages[0].Values["value"].(string))
	Pool.XAck(context.Background(), key, group, res.Val()[0].Messages[0].ID)
	return r, nil
}
