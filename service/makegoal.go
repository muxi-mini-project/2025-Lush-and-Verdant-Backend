package service

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
)

func (gsr *GoalService) PostGoal(message model.TasksData) error {
	for _, task := range message.Tasks {
		// 创建任务并关联到userID
		taskData := model.Task{
			Name:        task.Name,
			Description: task.Description,
			StartTime:   task.StartTime,
			EndTime:     task.EndTime,
			UserID:      task.UserID, // 关联任务与用户
			IsCompleted: false,
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
				return fmt.Errorf("事件保存失败")
			}
		}
	}
	return nil
}

func (gsr *GoalService) UpdateGoal(userID uint, message model.TasksData) error {
	for _, task := range message.Tasks {
		var existingTask model.Task
		if err := gsr.db.Where("id = ? AND user_id = ?", task.ID, userID).First(&existingTask).Error; err != nil {
			return fmt.Errorf("未找到该任务")
		}

		existingTask.Name = task.Name
		existingTask.Description = task.Description
		existingTask.StartTime = task.StartTime
		existingTask.EndTime = task.EndTime
		existingTask.IsCompleted = task.IsCompleted

		if err := gsr.db.Save(&existingTask).Error; err != nil {
			return fmt.Errorf("任务更新失败")
		}

		for _, event := range task.Events {
			var existingEvent model.Event
			if err := gsr.db.Where("id = ? AND task_id = ?", event.ID, event.TaskID).First(&existingEvent).Error; err != nil {
				return fmt.Errorf("未找到该事件")
			}

			existingEvent.Name = event.Name
			existingEvent.Description = event.Description
			existingEvent.StartTime = event.StartTime
			existingEvent.EndTime = event.EndTime

			if err := gsr.db.Save(&existingEvent).Error; err != nil {
				return fmt.Errorf("事件更新失败")
			}
		}
	}
	return nil
}

func (gsr *GoalService) HistoricalGoal(userID uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := gsr.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("获取目标失败")
	}

	for i := range tasks {
		var events []model.Event
		if err := gsr.db.Where("task_id = ?", tasks[i].ID).Find(&events).Error; err != nil {
			return nil, fmt.Errorf("获取事件失败")
		}
		tasks[i].Events = events
	}

	return tasks, nil
}

func (gsr *GoalService) DeleteGoal(userID uint, taskID string) error {
	var task model.Task

	if err := gsr.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return fmt.Errorf("未找到该任务")
	}

	if err := gsr.db.Where(&task).Error; err != nil {
		return fmt.Errorf("任务删除失败")
	}

	return nil
}
