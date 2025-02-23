package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Description  string `gorm:"type:text"`
	Password     string `gorm:"type:varchar(100);"`
	IsPublic     bool   `gorm:"default:true"`
	GroupOwnerId uint   `gorm:"unique"`
	Images       []GroupImage
	User         User `gorm:"foreignKey:GroupOwnerId"` // 关联到 User

}
