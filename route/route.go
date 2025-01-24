package route

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	r  *gin.Engine
	us *UserSvc
	ss *SloganSvc
	gs *GoalSvc
}

func NewApp(us *UserSvc, ss *SloganSvc, gs *GoalSvc) *App {
	r := gin.Default()
	us.NewUserGroup(r)
	ss.SloganGroup(r)
	gs.GoalGroup(r)

	return &App{
		r:  r,
		us: us,
		ss: ss,
		gs: gs,
	}
}

func (a *App) Run() {
	a.r.Run()
}
