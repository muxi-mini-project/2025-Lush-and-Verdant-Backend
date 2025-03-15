package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"strconv"
)

type GoalService interface {
	PostGoal(userID uint, req request.PostGoalRequest) (*model.Goal, error)
	UpdateTask(userID uint, taskID uint, req request.TaskRequest) error
	HistoricalGoal(userID uint) (map[string][]response.TaskWithChecks, error)
	DeleteTask(userID uint, taskID uint) error
	CheckTask(userID uint, taskID uint) (int, error)
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
func (gsr *GoalServiceImpl) PostGoal(userID uint, req request.PostGoalRequest) (*model.Goal, error) {
	newGoal := model.Goal{
		UserID: userID,
		Date:   req.Date,
	}
	if err := gsr.GoalDao.CreateGoal(&newGoal); err != nil {
		return nil, err
	}

	// 批量创建任务
	for _, task := range req.Tasks {
		newTask := model.Task{
			GoalID:  newGoal.ID,
			Title:   task.Title,
			Details: task.Details,
		}
		if err := gsr.GoalDao.CreateTask(&newTask); err != nil {
			return nil, err
		}
	}

	// 重新获取Goal，包含预加载的Tasks
	createdGoal, err := gsr.GoalDao.GetGoal(newGoal.ID, userID)
	if err != nil {
		return nil, err
	}

	return createdGoal, nil
}

// UpdateGoal 更新已有目标及任务
func (gsr *GoalServiceImpl) UpdateTask(userID uint, taskID uint, req request.TaskRequest) error {
	task, err := gsr.GoalDao.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("获取任务失败:%v", err)
	}

	// 验证任务所属权
	goal, err := gsr.GoalDao.GetGoal(task.GoalID, userID)
	if err != nil || goal.ID == 0 {
		return fmt.Errorf("无权操作该任务")
	}

	// 更新任务信息
	task.Title = req.Title
	task.Details = req.Details
	if err := gsr.GoalDao.UpdateTask(task); err != nil {
		return fmt.Errorf("更新失败:%v", err)
	}

	return nil
}

// HistoricalGoal 获取用户所有历史目标及任务完成情况
func (gsr *GoalServiceImpl) HistoricalGoal(userID uint) (map[string][]response.TaskWithChecks, error) {
	goals, err := gsr.GoalDao.GetGoals(userID)
	if err != nil {
		return nil, err
	}

	goalMap := make(map[string][]response.TaskWithChecks)
	for _, goal := range goals {
		tasks := make([]response.TaskWithChecks, 0)
		for _, task := range goal.Tasks {
			tasks = append(tasks, response.TaskWithChecks{
				TaskID:    strconv.Itoa(int(task.ID)),
				Title:     task.Title,
				Details:   task.Details,
				Completed: task.Completed,
			})
		}
		goalMap[goal.Date] = tasks
	}

	return goalMap, nil
}

// DeleteGoal 删除目标及任务
func (gsr *GoalServiceImpl) DeleteTask(userID uint, taskID uint) error {
	task, err := gsr.GoalDao.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("获取任务失败:%v", err)
	}

	// 验证任务所属权
	goal, err := gsr.GoalDao.GetGoal(task.GoalID, userID)
	if err != nil || goal.ID == 0 {
		return fmt.Errorf("无权操作该任务")
	}

	if err := gsr.GoalDao.DeleteTaskByID(taskID); err != nil {
		return fmt.Errorf("删除失败:%v", err)
	}

	return nil
}

// 标记任务为完成
func (gsr *GoalServiceImpl) CheckTask(userID uint, taskID uint) (int, error) {
	task, err := gsr.GoalDao.GetTaskByID(taskID)
	if err != nil {
		return 0, fmt.Errorf("获取任务失败:%v", err)
	}

	// 验证任务是否属于该用户
	goal, err := gsr.GoalDao.GetGoal(task.GoalID, userID)
	if err != nil || goal.ID == 0 {
		return 0, fmt.Errorf("无法访问该任务")
	}

	// 更新任务状态为已完成
	task.Completed = true
	if err := gsr.GoalDao.UpdateTask(task); err != nil {
		return 0, fmt.Errorf("更新任务状态失败:%v", err)
	}

	// 统计该目标下的已完成任务数量
	completedCount, err := gsr.GoalDao.CountCompletedTaskByGoal(goal.ID)
	if err != nil {
		return 0, fmt.Errorf("统计失败:%v", err)
	}

	return completedCount, nil
}
