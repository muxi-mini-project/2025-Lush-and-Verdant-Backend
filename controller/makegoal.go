package controller

import (
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	result := mc.gpt.AskForGoal(c)

	dataJSON, err := json.Marshal(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: "解析失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "请求成功", Data: string(dataJSON)})
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
