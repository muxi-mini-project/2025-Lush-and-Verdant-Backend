package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type GoalSvc struct {
	gc *controller.GoalController
}

func NewGoalSvc(gc *controller.GoalController) *GoalSvc {
	return &GoalSvc{
		gc: gc,
	}
}

func (g *GoalSvc) GoalGroup(r *gin.Engine) {
	r.Use(middleware.Cors())
	Goal := r.Group("/goal")
	{
		Goal.GET("/GetGoal", g.gc.GetGoal)
		Goal.POST("/MakeGoal", g.gc.PostGoal)
		Goal.PUT("/UpdateGoal", g.gc.UpdateGoal)
		Goal.GET("/HistoricalGoal", g.gc.HistoricalGoal)
		Goal.DELETE("/DeleteGoal/:task_id", g.gc.DeleteGoal)
	}
}
