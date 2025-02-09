package service

import (
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
)

type GoalService interface {
	PostGoal(model.TasksData) error
	UpdateGoal(uint, model.TasksData) error
	HistoricalGoal(uint) ([]*model.Task, error)
	DeleteGoal(uint, string) error
}

type GoalServiceImpl struct {
	GoalDao dao.GoalDAO
}

func NewGoalServiceImpl(goalDao dao.GoalDAO) *GoalServiceImpl {
	return &GoalServiceImpl{
		GoalDao: goalDao,
	}
}

func (gsr *GoalServiceImpl) PostGoal(message model.TasksData) error {
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

		err := gsr.GoalDao.CreatTask(&taskData)
		if err != nil {
			return err
		}

		for _, event := range task.Events {
			eventData := model.Event{
				Name:        event.Name,
				Description: event.Description,
				StartTime:   event.StartTime,
				EndTime:     event.EndTime,
				TaskID:      event.TaskID, // 关联任务
			}
			err := gsr.GoalDao.CreatEvent(&eventData)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gsr *GoalServiceImpl) UpdateGoal(userID uint, message model.TasksData) error {
	for _, task := range message.Tasks {
		existingTask, err := gsr.GoalDao.GetTask(task.ID, userID)
		if err != nil {
			return err
		}
		existingTask.Name = task.Name
		existingTask.Description = task.Description
		existingTask.StartTime = task.StartTime
		existingTask.EndTime = task.EndTime
		existingTask.IsCompleted = task.IsCompleted
		err = gsr.GoalDao.UpdateTask(existingTask)
		if err != nil {
			return err
		}
		for _, event := range task.Events {
			existingEvent, err := gsr.GoalDao.GetEvent(event.ID, event.TaskID)
			if err != nil {
				return err
			}
			existingEvent.Name = event.Name
			existingEvent.Description = event.Description
			existingEvent.StartTime = event.StartTime
			existingEvent.EndTime = event.EndTime

			err = gsr.GoalDao.UpdateEvent(existingEvent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gsr *GoalServiceImpl) HistoricalGoal(userID uint) ([]*model.Task, error) {
	tasks, err := gsr.GoalDao.GetTasks(userID)
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		events, err := gsr.GoalDao.GetEvents(tasks[i].ID)
		if err != nil {
			return nil, err
		}
		tasks[i].Events = events
	}

	return tasks, nil
}

func (gsr *GoalServiceImpl) DeleteGoal(userID uint, taskID string) error {
	err := gsr.GoalDao.DeleteTask(taskID, userID)
	if err != nil {
		return err
	}
	return nil
}
