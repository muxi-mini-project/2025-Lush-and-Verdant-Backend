package route

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	r  *gin.Engine
	us *UserSvc
	ss *SloganSvc
	gs *GoalSvc
	is *ImageSvc
}

func NewApp(us *UserSvc, ss *SloganSvc, gs *GoalSvc, is *ImageSvc) *App {
	r := gin.Default()
	us.NewUserGroup(r)
	ss.SloganGroup(r)
	gs.GoalGroup(r)
	is.ImageGroup(r)

	return &App{
		r:  r,
		us: us,
		ss: ss,
		gs: gs,
		is: is,
	}
}

func (a *App) Run() {
	a.r.Run()
}
