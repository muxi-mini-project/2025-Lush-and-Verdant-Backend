package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `gorm:"varchar(100);"`
	Password   string `gorm:"varchar(100);"`
	Email      string `gorm:"varchar(100);unique;default null;"`
	DeviceNum  string `gorm:"varchar(100);unique;default null;"`
	GoalPublic bool   `gorm:"default false;"`
	Slogan     string `gorm:"varchar(100);default null;"`
}
type Email struct {
	gorm.Model
	Email  string `gorm:"varchar(100);unique;default null;"`
	Code   string `gorm:"varchar(100);default null;"`
	Status bool   `gorm:"default false;"`
}
