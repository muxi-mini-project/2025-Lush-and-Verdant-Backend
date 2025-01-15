package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"github.com/gin-gonic/gin"
)

type UserSvc struct {
	uc *controller.UserController
}

func NewUserSvc(uc *controller.UserController) *UserSvc {
	return &UserSvc{
		uc: uc,
	}
}
func (u *UserSvc) NewUserGroup(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/send_email", u.uc.SendEmail)
		userGroup.POST("/register", u.uc.Register)
		userGroup.POST("/login", u.uc.Login)
		userGroup.POST("login_v", u.uc.Login_v)
		userGroup.POST("/foralt", u.uc.ForAlt)
		userGroup.POST("/cancel", u.uc.Cancel)

	}
}
