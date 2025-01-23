package controller

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

var timelayout = config.TimeLayout

type GoalController struct{}

func NewGoalController() *GoalController {
	return &GoalController{}
}

func (mc *GoalController) GetGoal(c *gin.Context) {
	result := tool.AskForSlogan(c)

	dataJSON, err := json.Marshal(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Error: "解析失败"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Message: "请求成功", Data: string(dataJSON)})
}

func (mc *GoalController) PostGoal(c *gin.Context) {
	var message model.TasksData
	if err := c.ShouldBind(message); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "解析失败"})
		return
	}

	db := dao.NewDB(dsn)

	for _, task := range message.Tasks {
		taskData := model.Task{
			Name:        task.Name,
			Description: task.Description,
			StartTime:   task.StartTime,
			EndTime:     task.EndTime,
		}

		if err := db.Create(&taskData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Error: "任务保存失败"})
			return
		}

		for _, event := range task.Events {
			eventData := model.Event{
				Name:        event.Name,
				Description: event.Description,
				StartTime:   event.StartTime,
				EndTime:     event.EndTime,
				TaskID:      event.TaskID, // 关联任务
			}

			if err := db.Create(&eventData).Error; err != nil {
				c.JSON(http.StatusInternalServerError, model.Response{Error: "事件保存失败"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, model.Response{Message: "保存成功"})
}
