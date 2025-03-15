package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type GoalDAO interface {
	CreateGoal(goal *model.Goal) error
	UpdateGoal(goal *model.Goal) error
	GetGoals(userId uint) ([]model.Goal, error)
	GetGoal(goalId uint, userId uint) (*model.Goal, error)
	DeleteGoal(goalId uint, userId uint) error

	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	GetTasks(goalId uint) ([]model.Task, error)
	DeleteTasks(goalId uint, UserId uint) error

	CountCompletedTaskByGoal(goalID uint) (int, error)
	GetTaskByID(taskID uint) (*model.Task, error)
	DeleteTaskByID(taskID uint) error
}

type GoalDAOImpl struct {
	db *gorm.DB
}

func NewGoalDAOImpl(db *gorm.DB) *GoalDAOImpl {
	return &GoalDAOImpl{
		db: db,
	}
}

// 创建目标
func (dao *GoalDAOImpl) CreateGoal(goal *model.Goal) error {
	if err := dao.db.Create(goal).Error; err != nil {
		return fmt.Errorf("目标创建失败:%v", err)
	}
	return nil
}

// 更新目标
func (dao *GoalDAOImpl) UpdateGoal(goal *model.Goal) error {
	if err := dao.db.Save(goal).Error; err != nil {
		return fmt.Errorf("更新目标失败:%v", err)
	}
	return nil
}

// 获取用户的所有目标及其任务
func (dao *GoalDAOImpl) GetGoals(userId uint) ([]model.Goal, error) {
	var goals []model.Goal
	if err := dao.db.Preload("Tasks").Where("user_id = ?", userId).Find(&goals).Error; err != nil {
		return nil, fmt.Errorf("获取目标失败:%v", err)
	}
	return goals, nil
}

// 根据ID获取某个目标
func (dao *GoalDAOImpl) GetGoal(goalId uint, userId uint) (*model.Goal, error) {
	var goal model.Goal
	if err := dao.db.Preload("Tasks").Where("id = ? AND user_id = ?", goalId, userId).First(&goal).Error; err != nil {
		return nil, fmt.Errorf("目标不存在")
	}
	return &goal, nil
}

// 删除目标及其关联的任务
func (dao *GoalDAOImpl) DeleteGoal(goalId uint, userId uint) error {
	if err := dao.db.Where("id = ? AND user_id = ?", goalId, userId).Delete(&model.Goal{}).Error; err != nil {
		return fmt.Errorf("目标删除失败:%v", err)
	}

	// 删除该目标下的所有任务
	if err := dao.DeleteTasks(goalId, userId); err != nil {
		return fmt.Errorf("任务删除失败:%v", err)
	}

	return nil
}

// 创建任务
func (dao *GoalDAOImpl) CreateTask(task *model.Task) error {
	if err := dao.db.Create(task).Error; err != nil {
		return fmt.Errorf("任务创建失败:%v", err)
	}
	return nil
}

// 更新任务
func (dao *GoalDAOImpl) UpdateTask(task *model.Task) error {
	if err := dao.db.Save(task).Error; err != nil {
		return fmt.Errorf("更新任务失败:%v", err)
	}
	return nil
}

// 获取某个Goal下的所有任务
func (dao *GoalDAOImpl) GetTasks(goalId uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := dao.db.Where("goal_id = ?", goalId).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("获取任务失败:%v", err)
	}
	return tasks, nil
}

// 删除某个Goal下的所有任务
func (dao *GoalDAOImpl) DeleteTasks(goalId uint, userId uint) error {
	if err := dao.db.Exec("DELETE t FROM tasks t INNER JOIN goals g ON t.goal_id = g.id WHERE g.id = ? AND g.user_id = ?", goalId, userId).Error; err != nil {
		return fmt.Errorf("任务删除失败:%v", err)
	}
	return nil
}

// 获取完成数量
func (dao *GoalDAOImpl) CountCompletedTaskByGoal(goalID uint) (int, error) {
	var count int64
	if err := dao.db.Model(&model.Task{}).Where("goal_id = ? AND completed = ?", goalID, true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计完成任务失败:%v", err)
	}
	return int(count), nil
}

// 根据ID获取任务
func (dao *GoalDAOImpl) GetTaskByID(taskID uint) (*model.Task, error) {
	var task model.Task
	if err := dao.db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, fmt.Errorf("任务不存在")
	}
	return &task, nil
}

// 根据ID删除任务
func (dao *GoalDAOImpl) DeleteTaskByID(taskID uint) error {
	if err := dao.db.Where("id = ?", taskID).Delete(&model.Task{}).Error; err != nil {
		return fmt.Errorf("删除任务失败: %v", err)
	}
	return nil
}
