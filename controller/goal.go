package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/service"
	"fmt"
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
// @Router /goal/GetGoal [put]
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
// @Success 200 {object} response.Response{data=response.PostGoalResponse} "保存成功"
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

	createdGoal, err := mc.gsr.PostGoal(userIDUint, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "保存失败"})
		return
	}

	// 提取Task IDs
	taskIDs := make([]string, len(createdGoal.Tasks))
	for i, task := range createdGoal.Tasks {
		taskIDs[i] = strconv.Itoa(int(task.ID))
	}

	responseData := response.PostGoalResponse{
		GoalID:  strconv.Itoa(int(createdGoal.ID)) + "10",
		TaskIDs: taskIDs,
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "保存成功", Data: responseData})
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 用户更新任务信息
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.PostGoalRequest true "更新的任务数据"
// @Success 200 {object} response.Response "任务更新成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/UpdateGoal/{task_id} [put]
func (mc *GoalController) UpdateTask(c *gin.Context) {
	var message request.TaskRequest // 绑定请求中的数据到message结构体
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Response{Code: 401, Message: "用户未授权"})
		return
	}

	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的任务ID"})
		fmt.Println(err)
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "类型转换失败"})
		return
	}
	userIDUint := uint(userIDInt)

	err = mc.gsr.UpdateTask(userIDUint, uint(taskID), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "任务更新失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "任务更新成功"})
}

// HistoricalGoal 查询历史目标
// @Summary 查询历史目标
// @Description 用户获取历史目标列表
// @Tags 目标管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.TaskWithChecks} "请求成功"
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

// UsersGoal 查询用户目标
// @Summary 查询用户目标
// @Description 用户获取目标列表
// @Tags 目标管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.TaskWithChecks} "请求成功"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/UsersGoal/{user_id} [get]
func (mc *GoalController) UsersGoal(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的用户ID"})
		fmt.Println(err)
		return
	}

	goals, err := mc.gsr.HistoricalGoal(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取历史目标失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "请求成功", Data: goals})
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 用户删除指定任务
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Success 200 {object} response.Response "任务删除成功"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/DeleteGoal/{task_id} [delete]
func (mc *GoalController) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的任务ID"})
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

	err = mc.gsr.DeleteTask(userIDUint, uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "任务删除失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "任务删除成功"})
}

// CheckTask 检查任务
// @Summary 检查任务
// @Description 用户检查指定任务
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Success 200 {object} response.Response{data=response.DailyCount} "任务检查成功"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/CheckTask/{task_id} [get]
func (mc *GoalController) CheckTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的任务ID"})
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

	dailyCount, err := mc.gsr.CheckTask(userIDUint, uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "任务检查失败"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "任务检查成功", Data: dailyCount})
}

// PostGoals 批量创建目标
// @Summary 批量创建目标
// @Description 用户批量创建新目标
// @Tags 目标管理
// @Accept json
// @Produce json
// @Param request body request.PostGoalRequests true "任务数据"
// @Success 200 {object} response.Response{data=response.PostGoalResponse} "保存成功"
// @Failure 400 {object} response.Response "解析失败"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/MakeGoals [post]
func (mc *GoalController) PostGoals(c *gin.Context) {
	var message request.PostGoalRequests
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "解析失败"})
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

	createdGoals, err := mc.gsr.PostGoals(userIDUint, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "保存失败"})
		return
	}

	var responseData []response.PostGoalResponse
	for _, goals := range createdGoals {
		taskIDs := make([]string, len(goals.Tasks))
		for i, task := range goals.Tasks {
			taskIDs[i] = strconv.Itoa(int(task.ID))
		}
		responseData = append(responseData, response.PostGoalResponse{
			GoalID:  strconv.Itoa(int(goals.ID)) + "10",
			TaskIDs: taskIDs,
		})
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "保存成功", Data: responseData})
}

// GetCompletedTasksCount 获取用户每日完成任务数
// @Summary 每日完成统计
// @Description 获取用户每个日期下已完成任务的数量
// @Tags 目标管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.CountResponse} "成功获取数据"
// @Failure 401 {object} response.Response "用户未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /goal/count/{user_id} [get]
func (mc *GoalController) GetCompletedTasksCount(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的用户ID"})
		fmt.Println(err)
		return
	}

	counts, err := mc.gsr.GetCompletedTaskCount(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取成功", Data: counts})
}
