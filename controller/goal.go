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

// GetGoal 获取目标
// @Summary 获取目标
// @Description 获取目标数据
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.Question true "请求问题"
// @Success 200 {object} response.Response "数据推送完成"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "推送数据失败"
// @Router /goal/GetGoal [get]
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

// PostGoal 创建目标
// @Summary 创建目标
// @Description 用户创建新目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body model.TasksData true "任务数据"
// @Success 200 {object} response.Response "保存成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/MakeGoal [post]
func (mc *GoalController) PostGoal(c *gin.Context) {
	var message model.TasksData
	if err := c.ShouldBind(message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "解析失败"})
		return
	}

	// 从上下文获取userID
	userID, exists := c.Get("user_id")
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

// UpdateGoal 更新目标
// @Summary 更新目标
// @Description 用户更新目标信息
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body model.TasksData true "更新的任务数据"
// @Success 200 {object} response.Response "目标更新成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/UpdateGoal [put]
func (mc *GoalController) UpdateGoal(c *gin.Context) {
	var message model.TasksData // 绑定请求中的数据到message结构体
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "解析失败"})
		return
	}

	userID, exists := c.Get("user_id")
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

// HistoricalGoal 查询历史目标
// @Summary 查询历史目标
// @Description 用户获取历史目标列表
// @Tags 目标管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "请求成功"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/HistoricalGoal [get]
func (mc *GoalController) HistoricalGoal(c *gin.Context) {
	userID, exists := c.Get("user_id")
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

// DeleteGoal 删除目标
// @Summary 删除目标
// @Description 用户删除指定目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Success 200 {object} response.Response "目标删除成功"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/DeleteGoal/{task_id} [delete]
func (mc *GoalController) DeleteGoal(c *gin.Context) {
	taskID := c.Param("task_id")

	userID, exists := c.Get("user_id")
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
