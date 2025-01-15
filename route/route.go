package route

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	r *gin.Engine
}

func NewApp(us *UserSvc) *App {
	r := gin.Default()
	us.NewUserGroup(r)
	return &App{
		r: r,
	}
}
func (a *App) Run() {
	a.r.Run()
}
