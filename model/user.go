package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `gorm:"type:varchar(100);"`
	Password   string `gorm:"type:varchar(100);"`
	Email      string `gorm:"type:varchar(100);unique;default:null;"`
	DeviceNum  string `gorm:"type:varchar(100);unique;default:null;"`
	GoalPublic bool   `gorm:"default:false;"`
	Slogan     string `gorm:"type:varchar(100);default:null;"`
	Images     []UserImage
	Groups     []Group `gorm:"many2many:user_groups;"` // 多对多关系
}
