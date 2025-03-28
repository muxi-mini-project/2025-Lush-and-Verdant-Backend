package model

import (
	"gorm.io/gorm"
	"time"
)

// Goal 结构体，包含日期和多个任务
type Goal struct {
	gorm.Model
	UserID uint   `gorm:"index:idx_user_date"` // 关联User
	Date   string `gorm:"index:idx_user_date"` // 日期
	Tasks  []Task `gorm:"foreignkey:GoalID"`   // 任务数组
}

// Task 结构体，表示具体任务
type Task struct {
	gorm.Model
	GoalID    uint   `gorm:"not null"` // 关联Goal
	Title     string `gorm:"varchar(255);not null"`
	Details   string `gorm:"text;not null"`
	Completed bool   `gorm:"default:false"`
}

type TaskCheck struct {
	gorm.Model
	TaskID    uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	CheckedAt time.Time `gorm:"not null"`
}
