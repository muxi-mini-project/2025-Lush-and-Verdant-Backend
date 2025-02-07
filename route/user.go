package route

import (
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func (u *UserSvc) NewUserGroup(r *gin.Engine) {
	r.Use(middleware.Cors()) //解决跨域问题
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
