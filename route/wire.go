package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"github.com/google/wire"
)

type UserSvc struct {
	uc *controller.UserController
}

func NewUserSvc(uc *controller.UserController) *UserSvc {
	return &UserSvc{
		uc: uc,
	}
}

type SloganSvc struct {
	uc *controller.SloganController
}

func NewSloganSvc(uc *controller.SloganController) *SloganSvc {
	return &SloganSvc{
		uc: uc,
	}
}

type GoalSvc struct {
	gc *controller.GoalController
}

func NewGoalSvc(gc *controller.GoalController) *GoalSvc {
	return &GoalSvc{
		gc: gc,
	}
}

var RouteSet = wire.NewSet(NewUserSvc, NewSloganSvc, NewGoalSvc, NewApp)
