package route

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	r   *gin.Engine
	us  *UserSvc
	ss  *SloganSvc
	gs  *GoalSvc
	is  *ImageSvc
	gsv *GroupSve
	cs  *ChatSve
}

func NewApp(us *UserSvc, ss *SloganSvc, gs *GoalSvc, is *ImageSvc, gsv *GroupSve, cs *ChatSve) *App {
	r := gin.Default()
	us.NewUserGroup(r)
	ss.SloganGroup(r)
	gs.GoalGroup(r)
	is.ImageGroup(r)
	gsv.Group(r)
	cs.ChatGroup(r)

	return &App{
		r:   r,
		us:  us,
		ss:  ss,
		gs:  gs,
		is:  is,
		gsv: gsv,
		cs:  cs,
	}
}

func (a *App) Run() {
	a.r.Run()
}
