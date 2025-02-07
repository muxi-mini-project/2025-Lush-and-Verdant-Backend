package controller

import (
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/google/wire"
)

type UserController struct {
	usr *service.UserService
}

func NewUserController(usr *service.UserService) *UserController {
	return &UserController{
		usr: usr,
	}
}

type SloganController struct {
	ssr *service.SloganService
}

func NewSloganController(ssr *service.SloganService) *SloganController {
	return &SloganController{
		ssr: ssr,
	}
}

type GoalController struct {
	gsr *service.GoalService
}

func NewGoalController(gsr *service.GoalService) *GoalController {
	return &GoalController{
		gsr: gsr,
	}
}

var ControllerSet = wire.NewSet(NewUserController, NewSloganController, NewGoalController)
