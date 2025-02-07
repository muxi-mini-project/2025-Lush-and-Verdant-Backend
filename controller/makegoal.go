package controller

import (
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/client"
	"2025-Lush-and-Verdant-Backend/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (mc *GoalController) GetGoal(c *gin.Context) {
	result := client.AskForGoal(c)

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
	err := mc.gsr.PostGoal(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "保存成功"})
}
