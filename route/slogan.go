package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type SloganSvc struct {
	uc  *controller.SloganController
	jwt *middleware.JwtClient
}

func NewSloganSvc(uc *controller.SloganController, jwt *middleware.JwtClient) *SloganSvc {
	return &SloganSvc{
		uc:  uc,
		jwt: jwt,
	}
}

func (u *SloganSvc) SloganGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	Slogans := r.Group("/slogan")
	{
		Slogans.GET("/GetSlogan/:device_num", u.uc.GetSlogan)
		Slogans.PUT("/ChangeSlogan/:user_id", u.uc.ChangeSlogan)

		Slogans.Use(u.jwt.AuthMiddleware())
		Slogans.GET("/SearchSlogan", u.uc.SearchSlogan)
	}
}
