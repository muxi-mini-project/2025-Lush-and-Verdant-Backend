package dao

import "gorm.io/gorm"

type GroupDAO interface {
}

type GroupDAOImpl struct {
	db *gorm.DB
}
