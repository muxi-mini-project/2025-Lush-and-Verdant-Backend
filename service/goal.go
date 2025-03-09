package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
)

type GoalService interface {
	PostGoal(userID uint, req request.PostGoalRequest) error
	UpdateGoal(userID uint, goalID uint, req request.PostGoalRequest) error
	HistoricalGoal(userID uint) (map[string][]model.Task, error)
	DeleteGoal(userID uint, goalID uint) error
}

type GoalServiceImpl struct {
	GoalDao dao.GoalDAO
}

func NewGoalServiceImpl(goalDao dao.GoalDAO) *GoalServiceImpl {
	return &GoalServiceImpl{
		GoalDao: goalDao,
	}
}

// PostGoal 创建新的目标及任务
func (gsr *GoalServiceImpl) PostGoal(userID uint, req request.PostGoalRequest) error {
	newGoal := model.Goal{
		UserID: userID,
		Date:   req.Date,
	}
	if err := gsr.GoalDao.CreateGoal(&newGoal); err != nil {
		return err
	}

	// 批量创建任务
	for _, task := range req.Tasks {
		newTask := model.Task{
			GoalID:  newGoal.ID,
			Title:   task.Title,
			Details: task.Details,
		}
		if err := gsr.GoalDao.CreateTask(&newTask); err != nil {
			return err
		}
	}

	return nil
}

// UpdateGoal 更新已有目标及任务
func (gsr *GoalServiceImpl) UpdateGoal(userID uint, goalID uint, req request.PostGoalRequest) error {
	goal, err := gsr.GoalDao.GetGoal(goalID, userID)
	if err != nil {
		return err
	}

	goal.Date = req.Date
	if err := gsr.GoalDao.CreateGoal(goal); err != nil {
		return err
	}

	// 先删除原有任务，再添加新任务
	if err := gsr.GoalDao.DeleteTasks(goal.ID); err != nil {
		return err
	}

	for _, task := range req.Tasks {
		newTask := model.Task{
			GoalID:  goal.ID,
			Title:   task.Title,
			Details: task.Details,
		}
		if err := gsr.GoalDao.CreateTask(&newTask); err != nil {
			return err
		}
	}

	return nil
}

// HistoricalGoal 获取用户所有历史目标及任务
func (gsr *GoalServiceImpl) HistoricalGoal(userID uint) (map[string][]model.Task, error) {
	goals, err := gsr.GoalDao.GetGoals(userID)
	if err != nil {
		return nil, err
	}

	goalMap := make(map[string][]model.Task)
	for _, goal := range goals {
		goalMap[goal.Date] = goal.Tasks
	}

	return goalMap, nil
}

// DeleteGoal 删除目标及任务
func (gsr *GoalServiceImpl) DeleteGoal(userID uint, goalID uint) error {
	return gsr.GoalDao.DeleteGoal(goalID, userID)
}
