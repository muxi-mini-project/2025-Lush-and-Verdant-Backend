package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type GoalDAO interface {
	CreateGoal(goal *model.Goal) error
	GetGoals(userId uint) ([]model.Goal, error)
	GetGoal(goalId uint, userId uint) (*model.Goal, error)
	DeleteGoal(goalId uint, userId uint) error

	CreateTask(task *model.Task) error
	GetTasks(goalId uint) ([]model.Task, error)
	DeleteTasks(goalId uint) error
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
	if err := dao.db.Where("goal_id = ?", goalId).Delete(&model.Task{}).Error; err != nil {
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

// 获取某个Goal下的所有任务
func (dao *GoalDAOImpl) GetTasks(goalId uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := dao.db.Where("goal_id = ?", goalId).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("获取任务失败:%v", err)
	}
	return tasks, nil
}

// 删除某个Goal下的所有任务
func (dao *GoalDAOImpl) DeleteTasks(goalId uint) error {
	if err := dao.db.Where("goal_id = ?", goalId).Delete(&model.Task{}).Error; err != nil {
		return fmt.Errorf("任务删除失败:%v", err)
	}
	return nil
}
