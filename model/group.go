package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Description  string `gorm:"type:text"`
	Password     string `gorm:"type:varchar(100);"`
	IsPublic     bool   `gorm:"default:true"`
	GroupOwnerId uint   // 确保与 User.ID 兼容
	// 外键字段
	GroupOwner User `gorm:"foreignKey:GroupOwnerId;constraint:OnDelete:CASCADE;"` // 关联到 User 表
	Images     []GroupImage
	Users      []User `gorm:"many2many:user_groups;"` // 多对多关系
}
