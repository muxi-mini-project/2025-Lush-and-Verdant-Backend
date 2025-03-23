package model

import (
	"gorm.io/gorm"
	"time"
)

type ForestLike struct {
	gorm.Model
	From      uint `gorm:"not null"`
	To        uint `gorm:"not null"`
	TimeStamp time.Time
}
