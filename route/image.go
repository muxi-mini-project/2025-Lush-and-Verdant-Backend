package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type ImageSvc struct {
	ic *controller.ImageController
}

func NewImageSvc(ic *controller.ImageController) *ImageSvc {
	return &ImageSvc{
		ic: ic,
	}
}

func (i *ImageSvc) ImageGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	r.GET("/get_token", i.ic.GetUpToken)
	userImage := r.Group("/image/user")
	{
		userImage.GET("/get/:id", i.ic.GetUserImage)
		userImage.GET("/history/:id", i.ic.GetUserAllImage)
		userImage.PUT("/update", i.ic.UpdateUserImage)
	}
	group := r.Group("/image/group")
	{
		group.GET("/:id", i.ic.GetGroupImage)
		group.GET("/history/:id", i.ic.GetGroupAllImage)
		group.PUT("/update", i.ic.UpdateGroupImage)
	}
}
