package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type LikeSvc struct {
	lc  *controller.LikeController
	jwt *middleware.JwtClient
}

func NewLikeSvc(lc *controller.LikeController, jwt *middleware.JwtClient) *LikeSvc {
	return &LikeSvc{
		lc:  lc,
		jwt: jwt,
	}
}

func (l *LikeSvc) LikeGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	LikeGroup := r.Group("/like")
	{
		LikeGroup.Use(l.jwt.AuthMiddleware())
		LikeGroup.POST("/send", l.lc.Like)
		LikeGroup.GET("/get/:to", l.lc.GetForestAllLikes)
		LikeGroup.POST("/status", l.lc.GetForestLikeStatus)
	}
}
