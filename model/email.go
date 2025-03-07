package model

import "gorm.io/gorm"

type Email struct {
	gorm.Model
	Email string `gorm:"varchar(100);unique;default:null;"`
}
