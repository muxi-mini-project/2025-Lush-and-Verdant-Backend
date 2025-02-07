package service

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

type SloganService struct {
	db *gorm.DB
}

func NewSloganService(db *gorm.DB) *SloganService {
	return &SloganService{
		db: db,
	}
}

type GoalService struct {
	db *gorm.DB
}

func NewGoalService(db *gorm.DB) *GoalService {
	return &GoalService{
		db: db,
	}
}

var ServiceSet = wire.NewSet(NewUserService, NewSloganService, NewGoalService)
