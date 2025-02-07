package route

import (
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func (g *GoalSvc) GoalGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	Goal := r.Group("/goal")
	{
		Goal.GET("/GetGoal", g.gc.GetGoal)
		Goal.POST("/MakeGoal", g.gc.PostGoal)
	}
}
