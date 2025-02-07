package service

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
)

func (gsr *GoalService) PostGoal(message model.TasksData) error {

	for _, task := range message.Tasks {
		taskData := model.Task{
			Name:        task.Name,
			Description: task.Description,
			StartTime:   task.StartTime,
			EndTime:     task.EndTime,
		}

		if err := gsr.db.Create(&taskData).Error; err != nil {
			return fmt.Errorf("任务保存失败")
		}

		for _, event := range task.Events {
			eventData := model.Event{
				Name:        event.Name,
				Description: event.Description,
				StartTime:   event.StartTime,
				EndTime:     event.EndTime,
				TaskID:      event.TaskID, // 关联任务
			}

			if err := gsr.db.Create(&eventData).Error; err != nil {
				return fmt.Errorf("时间保存失败")
			}
		}
	}
	return nil
}
