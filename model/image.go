package model

import "gorm.io/gorm"

type UserImage struct {
	gorm.Model
	Url    string `gorm:"type:varchar(255)"`
	UserID uint
	User   User `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE" `
}

type GroupImage struct {
	gorm.Model
	Url     string `gorm:"type:varchar(255)"`
	GroupID uint
	Group   Group `gorm:"foreignkey:GroupID;constraint:OnDelete:CASCADE"`
}
