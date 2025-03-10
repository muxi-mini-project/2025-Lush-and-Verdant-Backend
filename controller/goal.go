package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// GetGoal 获取AI生成的目标
// @Summary 获取目标
// @Description 获取目标数据
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.Question true "请求问题"
// @Success 200 {object} response.Response{data=response.Goals} "数据推送完成"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "推送数据失败"
// @Router /goal/GetGoal [get]
func (mc *GoalController) GetGoal(c *gin.Context) {
	var message request.Question
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "请求参数错误"})
		return
	}

	// 将数据传递给service层，获取返回结果
	result := mc.gpt.AskForGoal(c, message)
	if result == nil {
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "AI获取目标成功", Data: result})
}

// PostGoal 创建目标
// @Summary 创建目标
// @Description 用户创建新目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.PostGoalRequest true "任务数据"
// @Success 200 {object} response.Response "保存成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/MakeGoal [post]
func (mc *GoalController) PostGoal(c *gin.Context) {
	var message request.PostGoalRequest
	if err := c.ShouldBind(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}

	// 从上下文获取userID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	err := mc.gsr.PostGoal(userIDUint, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "保存失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "保存成功"})
}

// UpdateGoal 更新目标
// @Summary 更新目标
// @Description 用户更新目标信息
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.PostGoalRequest true "更新的任务数据"
// @Success 200 {object} response.Response "目标更新成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/UpdateGoal/{goal_id} [put]
func (mc *GoalController) UpdateGoal(c *gin.Context) {
	var message request.PostGoalRequest // 绑定请求中的数据到message结构体
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	goalIDStr := c.Param("goal_id")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的目标ID"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	err = mc.gsr.UpdateGoal(userIDUint, uint(goalID), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "目标更新失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "目标更新成功"})
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
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	goals, err := mc.gsr.HistoricalGoal(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取历史目标失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "请求成功", Data: goals})
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
// @Router /goal/DeleteGoal/{goal_id} [delete]
func (mc *GoalController) DeleteGoal(c *gin.Context) {
	goalIDStr := c.Param("goal_id")

	goalID, err := strconv.ParseUint(goalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的目标ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	err = mc.gsr.DeleteGoal(userIDUint, uint(goalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "目标删除失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "目标删除成功"})
}
