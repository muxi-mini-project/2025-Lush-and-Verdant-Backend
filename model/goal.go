package model

import (
	"time"
)

type CustomTime time.Time

// Event 结构体表示单个事件
type Event struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartTime   CustomTime `json:"start_time"`
	EndTime     CustomTime `json:"end_time"`
	TaskID      uint       `json:"task_id"` // 外键，关联任务
}

// Task 结构体表示任务，包含多个事件
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartTime   CustomTime `json:"start_time"`
	EndTime     CustomTime `json:"end_time"`
	UserID      uint       `json:"user_id"`      // 外键，关联用户
	Events      []Event    `json:"events"`       // 任务下的多个事件
	IsCompleted bool       `json:"is_completed"` // 检测任务是否完成
}

// TasksData 是最外层结构体，包含多个任务
type TasksData struct {
	Tasks []Task `json:"tasks"` // 一个任务列表
}
