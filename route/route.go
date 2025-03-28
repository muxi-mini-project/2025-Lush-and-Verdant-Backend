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
	ls  *LikeSvc
}

func NewApp(us *UserSvc, ss *SloganSvc, gs *GoalSvc, is *ImageSvc, gsv *GroupSve, cs *ChatSve, ls *LikeSvc) *App {
	r := gin.Default()
	us.NewUserGroup(r)
	ss.SloganGroup(r)
	gs.GoalGroup(r)
	is.ImageGroup(r)
	gsv.Group(r)
	cs.ChatGroup(r)
	ls.LikeGroup(r)

	return &App{
		r:   r,
		us:  us,
		ss:  ss,
		gs:  gs,
		is:  is,
		gsv: gsv,
		cs:  cs,
		ls:  ls,
	}
}

func (a *App) Run() {
	a.r.Run()
}
