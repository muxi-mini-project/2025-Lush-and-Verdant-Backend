package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type GroupSve struct {
	gc *controller.GroupController
}

func NewGroupSve(gc *controller.GroupController) *GroupSve {
	return &GroupSve{
		gc: gc,
	}
}

func (g *GroupSve) Group(r *gin.Engine) {
	r.Use(middleware.Cors())
	group := r.Group("/group")
	{
		group.POST("/create", g.gc.CreateGroup)
		group.POST("/update", g.gc.UpdateGroup)
		group.POST("/delete", g.gc.DeleteGroup)
		group.GET("/info/:groupNum", g.gc.GetGroupInfo)
		group.GET("/members/:groupNum", g.gc.GetGroupMemberList)
		group.GET("/list/:id", g.gc.GetGroupList)
		group.POST("/member/add", g.gc.AddGroupMember)
		group.POST("/member/delete", g.gc.DeleteGroupMember)
		group.GET("/ten", g.gc.GetTenGroup)
		group.POST("/check", g.gc.CheckGroupMember)
		group.GET("/find", g.gc.FindGroup)
	}
}
