package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type GoalDAO interface {
	CreatTask(taskData *model.Task) error
	CreatEvent(eventData *model.Event) error
	GetTask(taskId uint, userId uint) (*model.Task, error)
	GetTasks(userId uint) ([]*model.Task, error)
	GetEvent(eventId uint, taskId uint) (*model.Event, error)
	GetEvents(taskId uint) ([]model.Event, error)
	UpdateTask(existingTask *model.Task) error
	UpdateEvent(existingEvent *model.Event) error
	DeleteTask(taskId string, userId uint) error
}
type GoalDAOImpl struct {
	db *gorm.DB
}

func NewGoalDAOImpl(db *gorm.DB) *GoalDAOImpl {
	return &GoalDAOImpl{
		db: db,
	}
}

// 创建任务
func (dao *GoalDAOImpl) CreatTask(taskData *model.Task) error {
	if err := dao.db.Create(&taskData).Error; err != nil {
		return fmt.Errorf("任务保存失败")
	}
	return nil
}

// 创建事件
func (dao *GoalDAOImpl) CreatEvent(eventData *model.Event) error {
	if err := dao.db.Create(&eventData).Error; err != nil {
		return fmt.Errorf("事件保存失败")
	}
	return nil
}

// 根据userId获取tasks
func (dao *GoalDAOImpl) GetTasks(userId uint) ([]*model.Task, error) {
	var tasks []*model.Task
	if err := dao.db.Where("user_id = ?", userId).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("获取目标失败")
	}
	return tasks, nil
}

// 根据task_id 和 user_id 查询任务
func (dao *GoalDAOImpl) GetTask(taskId uint, userId uint) (*model.Task, error) {
	var existingTask model.Task
	if err := dao.db.Where("id = ? AND user_id = ?", taskId, userId).First(&existingTask).Error; err != nil {
		return nil, fmt.Errorf("未找到该任务")
	}
	return &existingTask, nil
}

// 根据event_id 和 even_taskID 来查找event事件
func (dao *GoalDAOImpl) GetEvent(eventId uint, taskId uint) (*model.Event, error) {
	var existingEvent model.Event
	if err := dao.db.Where("id = ? AND task_id = ?", eventId, taskId).First(&existingEvent).Error; err != nil {
		return nil, fmt.Errorf("未找到该事件")
	}
	return &existingEvent, nil
}

// 根据taskId查询所有事件 //todo 这里没有用指针类型来传输数据
func (dao *GoalDAOImpl) GetEvents(taskId uint) ([]model.Event, error) {
	var events []model.Event
	if err := dao.db.Where("task_id = ?", taskId).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("获取事件失败")
	}
	return events, nil
}

// 更新任务
func (dao *GoalDAOImpl) UpdateTask(existingTask *model.Task) error {
	if err := dao.db.Save(&existingTask).Error; err != nil {
		return fmt.Errorf("任务更新失败")
	}
	return nil
}

// 更新事件
func (dao *GoalDAOImpl) UpdateEvent(existingEvent *model.Event) error {
	if err := dao.db.Save(&existingEvent).Error; err != nil {
		return fmt.Errorf("事件更新失败")
	}
	return nil
}

// 根据taskId和userId删除task
func (dao *GoalDAOImpl) DeleteTask(taskId string, userId uint) error {
	var task model.Task

	if err := dao.db.Where("id = ? AND user_id = ?", taskId, userId).First(&task).Error; err != nil {
		return fmt.Errorf("未找到该任务")
	}

	if err := dao.db.Where(&task).Error; err != nil {
		return fmt.Errorf("任务删除失败")
	}
	return nil
}
