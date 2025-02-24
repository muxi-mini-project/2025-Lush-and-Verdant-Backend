package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type GoalController struct {
	gsr service.GoalService
	gpt *client.ChatGptClient
}

func NewGoalController(gsr service.GoalService, gpt *client.ChatGptClient) *GoalController {
	return &GoalController{
		gsr: gsr,
		gpt: gpt,
	}
}

func (mc *GoalController) GetGoal(c *gin.Context) {
	var message request.Question
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "请求参数错误"})
		return
	}

	// 将数据传递给service层，获取返回结果
	result := mc.gpt.AskForGoal(c, message)

	// 初始化推送状态
	pushed := false

	// 使用SSE流式推送数据
	c.Stream(func(w io.Writer) bool {
		for key, value := range result {
			// 构建要发送的消息消息
			msg := fmt.Sprintf(`{"%s","%s"}`, key, value)

			// 将消息写入响应流
			_, err := c.Writer.Write([]byte("data:" + msg + "\n\n"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.Response{Error: "推送数据失败"})
				return false
			}

			// 强制推送消息到客户端
			c.Writer.Flush()

			// 标记已推送
			pushed = true
		}
		return true
	})

	if pushed {
		c.JSON(http.StatusOK, response.Response{Message: "数据推送完成"})
	} else {
		c.JSON(http.StatusInternalServerError, response.Response{Error: "无数据被推送"})
	}
}

func (mc *GoalController) PostGoal(c *gin.Context) {
	var message model.TasksData
	if err := c.ShouldBind(message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "解析失败"})
		return
	}

	// 从上下文获取userID
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Error: "用户未授权"})
		return
	}

	// 将userID添加到每个任务中
	for i := range message.Tasks {
		message.Tasks[i].UserID = userID.(uint)
	}

	err := mc.gsr.PostGoal(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "保存成功"})
}

func (mc *GoalController) UpdateGoal(c *gin.Context) {
	var message model.TasksData // 绑定请求中的数据到message结构体
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "解析失败"})
		return
	}

	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Error: "用户未授权"})
		return
	}

	for i := range message.Tasks {
		message.Tasks[i].UserID = userID.(uint) // 将userID关联到任务中
	}

	err := mc.gsr.UpdateGoal(userID.(uint), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "目标更新成功"})
}

func (mc *GoalController) HistoricalGoal(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Error: "用户未授权"})
		return
	}

	goals, err := mc.gsr.HistoricalGoal(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "请求成功", Data: goals})
}

func (mc *GoalController) DeleteGoal(c *gin.Context) {
	taskID := c.Param("task_id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Error: "用户未授权"})
		return
	}

	err := mc.gsr.DeleteGoal(userID.(uint), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "目标删除成功"})
}
