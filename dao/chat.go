package dao

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type ChatDAO interface {
	ReadMessage(stream, lastId string) ([]redis.XMessage, error)
	GetUserLastId(userId string) string
	GetGroupLastId(userId, groupId string) string
	GetUserHistory(stream string) []redis.XMessage
	GetGroupHistory(stream string) []redis.XMessage
	SetUserLastId(userId, lastId string)
	SetGroupLastId(userId, groupId, groupLastId string)
	AddUserMessage(personalStream string, message *request.Message) (string, error)
	AddGroupMessage(groupStream string, message *request.Message) (string, error)
	AddUserHistoryMessage(personalHistoryStream string, message *request.Message)
}

type ChatDAOImpl struct {
	rdb *redis.Client
}

func NewChatDAOImpl(rdb *redis.Client) *ChatDAOImpl {
	return &ChatDAOImpl{
		rdb: rdb,
	}
}

// ReadMessage 读取消息
func (dao *ChatDAOImpl) ReadMessage(stream, lastId string) ([]redis.XMessage, error) {
	message, err := dao.rdb.XRange(context.TODO(), stream, lastId, "+").Result()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// GetUserLastId 获取用户上次读取的时候
func (dao *ChatDAOImpl) GetUserLastId(userId string) string {
	lastId, _ := dao.rdb.HGet(context.TODO(), "user:last_read", userId).Result()

	if lastId == "" {
		lastId = "0"
	}
	return lastId
}

func (dao *ChatDAOImpl) GetGroupLastId(userId, groupId string) string {
	groupLastId, _ := dao.rdb.HGet(context.TODO(), "user:"+userId+":group_last_read", groupId).Result()
	if groupLastId == "" {
		groupLastId = "0"
	}
	return groupLastId
}

// SetUserLastId 更新读取的位置
func (dao *ChatDAOImpl) SetUserLastId(userId, lastId string) {
	dao.rdb.HSet(context.TODO(), "user:last_read", userId, lastId)
}

func (dao *ChatDAOImpl) SetGroupLastId(userId, groupId, groupLastId string) {
	dao.rdb.HSet(context.TODO(), "user:"+userId+":group_last_read", groupId, groupLastId)
}

// AddUserMessage 写入接收者的消息流
func (dao *ChatDAOImpl) AddUserMessage(personalStream string, message *request.Message) (string, error) {
	msgId, err := dao.rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: personalStream,
		Values: map[string]interface{}{
			"from":    message.From,
			"to":      message.To,
			"content": message.Content,
			"type":    "personal",
			"time":    time.Now().Format(time.DateTime),
		},
	}).Result()
	return msgId, err
}

func (dao *ChatDAOImpl) AddUserHistoryMessage(personalHistoryStream string, message *request.Message) {
	_, err := dao.rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: personalHistoryStream,
		Values: map[string]interface{}{
			"from":    message.From,
			"to":      message.To,
			"content": message.Content,
			"type":    "personal",
			"time":    time.Now().Format(time.DateTime),
		},
	}).Result()
	if err != nil {
		log.Printf("%s->%s添加历史消息失败", message.From, message.To)
		return
	}
}
func (dao *ChatDAOImpl) AddGroupMessage(groupStream string, message *request.Message) (string, error) {
	msgId, err := dao.rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: groupStream,
		Values: map[string]interface{}{
			"from":    message.From,
			"to":      message.To,
			"content": message.Content,
			"type":    "group",
			"time":    time.Now().Format(time.DateTime),
		},
	}).Result()
	return msgId, err
}

func (dao *ChatDAOImpl) GetUserHistory(stream string) []redis.XMessage {
	// 获取用户所有的消息流
	messages, err := dao.rdb.XRange(context.TODO(), stream, "-", "+").Result()
	if err != nil {
		log.Printf("获取用户\"%s\"的历史消息失败", stream)
		return nil
	}
	return messages

}

func (dao *ChatDAOImpl) GetGroupHistory(stream string) []redis.XMessage {
	// 获取用户所有的消息流
	messages, err := dao.rdb.XRange(context.TODO(), stream, "-", "+").Result()
	if err != nil {
		log.Printf("获取群聊\"%s\"的历史消息失败", stream)
		return nil
	}
	return messages

}
