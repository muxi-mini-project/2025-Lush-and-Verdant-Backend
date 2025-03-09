package model

import "gorm.io/gorm"

// Goal 结构体，包含日期和多个任务
type Goal struct {
	gorm.Model
	UserID uint   `gorm:"not null"`             // 关联User
	Date   string `gorm:"varchar(20);not null"` // 日期
	Tasks  []Task `gorm:"foreignkey:GoalID"`    // 任务数组
}

// Task 结构体，表示具体任务
type Task struct {
	gorm.Model
	GoalID  uint   `gorm:"not null"` // 关联Goal
	Title   string `gorm:"varchar(255);not null"`
	Details string `gorm:"text;not null"`
}
