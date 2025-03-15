package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChatController struct {
	csr service.ChatService
}

func NewChatController(csr service.ChatService) *ChatController {
	return &ChatController{csr: csr}
}

// HandleWebSocket 处理WebSocket连接
// @Summary 处理WebSocket连接
// @Description 通过身份验证获取用户ID并处理WebSocket连接
// @Tags Chat
// @Accept json
// @Produce json
// @Param userID header int true "用户ID"
// @Success 200 {object} response.Response "成功处理WebSocket连接"
// @Failure 400 {object} response.Response "用户ID不是整数"
// @Router /chat/ws [get]
func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	// 通过身份验证获取
	id, _ := c.Get("userID")
	userId, ok := id.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "id must be int",
		})
		return // 添加return以避免继续执行
	}
	idStr := strconv.Itoa(userId)
	cc.csr.HandleWebSocket(c.Writer, c.Request, idStr)
}

// GetUserHistory 获取用户历史消息
// @Summary 获取用户历史消息
// @Description 根据用户ID获取用户之间的历史消息
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body request.UserHistory true "请求参数"
// @Success 200 {object} response.Response{data=response.Messages} "成功返回用户历史消息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /chat/user/history [post]
func (cc *ChatController) GetUserHistory(c *gin.Context) {
	var req request.UserHistory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "获取用户历史消息失败"})
		return
	}
	messages := cc.csr.GetUserHistory(req.From, req.To)
	messagesRes := response.Messages{Messages: messages}
	c.JSON(http.StatusOK, response.Response{Code: 200, Data: messagesRes})
}

// GetGroupHistory 获取群组历史消息
// @Summary 获取群组历史消息
// @Description 根据群组ID获取群组的历史消息
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body request.GroupHistory true "请求参数"
// @Success 200 {object} response.Response{data=response.Messages} "成功返回群组历史消息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /chat/group/history [post]
func (cc *ChatController) GetGroupHistory(c *gin.Context) {
	var req request.GroupHistory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "获取群组历史消息失败"})
		return
	}
	messages := cc.csr.GetGroupHistory(req.GroupId)
	messagesRes := response.Messages{Messages: messages}
	c.JSON(http.StatusOK, response.Response{Code: 200, Data: messagesRes})
}
