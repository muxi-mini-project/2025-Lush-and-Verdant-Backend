package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/tool"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ChatService interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request, id string)
	GetUserHistory(from, to string) []response.Message
	GetGroupHistory(groupId string) []response.Message
}

type ChatServiceImpl struct {
	ctx         context.Context
	connections *sync.Map // 存储在线用户的连接
	upgrader    websocket.Upgrader
	Dao         dao.ChatDAO
	GroupDao    dao.GroupDAO
}

func NewChatServiceImpl(Dao dao.ChatDAO, groupDAO dao.GroupDAO) *ChatServiceImpl {
	return &ChatServiceImpl{
		ctx:         context.Background(),
		connections: &sync.Map{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Dao:      Dao,
		GroupDao: groupDAO,
	}

}

// HandleWebSocket 处理 WebSocket 连接
func (csr *ChatServiceImpl) HandleWebSocket(w http.ResponseWriter, r *http.Request, id string) {
	// 升级为 WebSocket 连接
	conn, err := csr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// 注册客户端
	csr.connections.Store(id, conn)

	go csr.SendHistoryMessage(id, conn)

	// 读取客户端消息
	for {
		var message request.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("WebSocket read message failed: %v", err)
			csr.connections.Delete(id)
			return
		}

		ok := csr.CheckMessage(&message)
		// 检测消息是否正确
		if ok {
			csr.HandleMessage(&message)
		} else {
			log.Printf("消息格式不正确")
		}
	}

}

// CheckMessage 检查消息格式是否正确
func (csr *ChatServiceImpl) CheckMessage(message *request.Message) bool {
	if message.Type != "group" && message.Type != "personal" {
		return false
	}
	if message.From == "" || message.To == "" || message.Content == "" {
		return false
	}
	return true
}

func (csr *ChatServiceImpl) HandleMessage(message *request.Message) {
	switch message.Type {
	case "personal":
		csr.HandlePersonalMessage(message)
	case "group":
		csr.HandleGroupMessage(message)
	}
}

func (csr *ChatServiceImpl) HandlePersonalMessage(message *request.Message) {
	//写入接收者的消息流
	personalStream := "user:msg:" + message.To
	//开一个协程写保存聊天历史
	ma, mi := tool.CompareString(message.From, message.To)
	personalHistoryStream := "user:msg:" + fmt.Sprintf("%d-%d", mi, ma)
	go csr.Dao.AddUserHistoryMessage(personalHistoryStream, message)

	//写入redis stream
	_, err := csr.Dao.AddUserMessage(personalStream, message)
	if err != nil {
		log.Printf("AddUserMessage failed: %v", err)
		return
	}

	//实时推送
	//查看接受者是否在线
	if conn, ok := csr.connections.Load(message.To); ok {
		response := map[string]interface{}{
			"type":    "personal",
			"from":    message.From,
			"to":      message.To,
			"content": message.Content,
			"time":    time.Now().Format(time.DateTime),
		}
		if err := conn.(*websocket.Conn).WriteJSON(response); err != nil {
			log.Printf("用户%s发送消息失败", message.From)
			csr.connections.Delete(message.To)
			return
		}
	}
}

func (csr *ChatServiceImpl) HandleGroupMessage(message *request.Message) {
	// 写入群聊消息流
	groupStream := "group:msg:" + message.To

	//写入redis stream
	_, err := csr.Dao.AddGroupMessage(groupStream, message)
	if err != nil {
		log.Printf("AddGroupMessage failed: %v", err)
		return
	}

	//实时推送
	//获取群聊的人员列表
	to, err := tool.StringToUint(message.To)
	if err != nil {
		log.Printf("StringToUint failed: %v", err)
		return
	}
	memberIDS, err := csr.GroupDao.GetGroupMemberIdList(to)
	if err != nil {
		log.Printf("GetGroupMemberIdList failed: %v", err)
		return
	}
	for _, memberID := range memberIDS {
		//查看接受者是否在线
		memberIdStr := strconv.Itoa(memberID)
		if conn, ok := csr.connections.Load(memberIdStr); ok {
			response := map[string]interface{}{
				"type":    "group",
				"from":    message.From,
				"to":      message.To,
				"content": message.Content,
				"time":    time.Now().Format(time.DateTime),
			}
			if err := conn.(*websocket.Conn).WriteJSON(response); err != nil {
				log.Printf("用户%s发送消息失败", message.From)
				csr.connections.Delete(memberIdStr)
				return
			}
		}
	}

}

// SendHistoryMessage 发送离线时候的消息
func (csr *ChatServiceImpl) SendHistoryMessage(userId string, conn *websocket.Conn) {
	//个人消息
	personalStream := "user:msg:" + userId
	//获取上次读取的消息的位置
	lastId := csr.Dao.GetUserLastId(userId)

	message, err := csr.Dao.ReadMessage(personalStream, lastId)
	if err != nil {
		log.Printf("Read message failed: %v", err)
		return
	}
	for _, v := range message {
		conn.WriteJSON(map[string]interface{}{
			"type":    "personal",
			"from":    v.Values["from"],
			"to":      v.Values["to"],
			"content": v.Values["content"],
			"time":    v.Values["time"],
		})
		//更新读取的位置
		lastId = v.ID
	}
	// 更新读取的位置
	csr.Dao.SetUserLastId(userId, lastId)

	//群聊消息
	//获取用户的群聊列表
	userIdUint, err := tool.StringToUint(userId)
	if err != nil {
		log.Printf("StringToUint failed: %v", err)
		return
	}
	groupIdS, err := csr.GroupDao.GetGroupIdList(userIdUint)
	if err != nil {
		log.Printf("GetGroupIdList failed: %v", err)
		return
	}
	for _, groupId := range groupIdS {
		//一个一个处理
		groupIdStr := strconv.Itoa(groupId)
		groupStream := "group:msg:" + groupIdStr
		groupLastId := csr.Dao.GetGroupLastId(userId, groupIdStr)

		groupMessages, err := csr.Dao.ReadMessage(groupStream, groupLastId)
		if err != nil {
			log.Printf("Read message failed: %v", err)
			return
		}
		for _, v := range groupMessages {
			conn.WriteJSON(map[string]interface{}{
				"type":    "group",
				"from":    v.Values["from"],
				"to":      v.Values["to"], //群聊id
				"content": v.Values["content"],
				"time":    v.Values["time"],
			})
			//更新读取的位置
			groupLastId = v.ID
		}
		// 更新读取的位置
		csr.Dao.SetGroupLastId(userId, groupIdStr, groupLastId)
	}
}

// GetUserHistory 获取用户历史聊天记录
func (csr *ChatServiceImpl) GetUserHistory(from, to string) []response.Message {
	ma, mi := tool.CompareString(from, to)
	stream := "user:msg:" + fmt.Sprintf("%d-%d", mi, ma)
	messagesHistory := csr.Dao.GetUserHistory(stream)
	var messages []response.Message
	for _, v := range messagesHistory {
		massage := response.Message{
			Type:    v.Values["type"].(string),
			From:    v.Values["from"].(string),
			To:      v.Values["to"].(string),
			Content: v.Values["content"].(string),
			Time:    v.Values["time"].(string),
		}
		messages = append(messages, massage)
	}
	return messages
}

// 获取群聊历史消息

func (csr *ChatServiceImpl) GetGroupHistory(groupId string) []response.Message {
	stream := "group:msg:" + groupId
	messagesHistory := csr.Dao.GetGroupHistory(stream)
	var messages []response.Message
	for _, v := range messagesHistory {
		massage := response.Message{
			Type:    v.Values["type"].(string),
			From:    v.Values["from"].(string),
			To:      v.Values["to"].(string),
			Content: v.Values["content"].(string),
			Time:    v.Values["time"].(string),
		}
		messages = append(messages, massage)
	}
	return messages
}

// todo 删除历史聊天记录
