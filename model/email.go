package model

import "gorm.io/gorm"

type Email struct {
	gorm.Model
	Email  string `gorm:"varchar(100);unique;default:null;"`
	Code   string `gorm:"varchar(100);default:null;"`
	Status bool   `gorm:"default:false;"`
}
