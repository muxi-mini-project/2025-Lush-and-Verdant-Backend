package route

import (
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func (u *SloganSvc) SloganGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	Slogans := r.Group("/slogan")
	{
		Slogans.GET("/GetSlogan/:device_num", u.uc.GetSlogan)
		Slogans.PUT("/ChangeSlogan", u.uc.ChangeSlogan)
	}
}
