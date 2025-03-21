package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type ChatSve struct {
	cc  *controller.ChatController
	jwt *middleware.JwtClient
}

func NewChatSve(cc *controller.ChatController, jwt *middleware.JwtClient) *ChatSve {
	return &ChatSve{cc: cc, jwt: jwt}
}

func (cs *ChatSve) ChatGroup(r *gin.Engine) {
	r.Use(middleware.Cors())

	Chat := r.Group("/chat")
	{
		Chat.GET("/ws", cs.cc.HandleWebSocket)
		Chat.Use(cs.jwt.AuthMiddleware())
		Chat.POST("/user/history", cs.cc.GetUserHistory)
		Chat.POST("/group/history", cs.cc.GetGroupHistory)
	}
}
