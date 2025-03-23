package route

import (
	"2025-Lush-and-Verdant-Backend/controller"
	"2025-Lush-and-Verdant-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type GoalSvc struct {
	gc  *controller.GoalController
	jwt *middleware.JwtClient
}

func NewGoalSvc(gc *controller.GoalController, jwt *middleware.JwtClient) *GoalSvc {
	return &GoalSvc{
		gc:  gc,
		jwt: jwt,
	}
}

func (g *GoalSvc) GoalGroup(r *gin.Engine) {
	r.Use(middleware.Cors())

	Goal := r.Group("/goal")
	{
		Goal.Use(g.jwt.AuthMiddleware())
		Goal.POST("/GetGoal", g.gc.GetGoal)
		Goal.POST("/MakeGoal", g.gc.PostGoal)
		Goal.PUT("/UpdateGoal/:task_id", g.gc.UpdateTask)
		Goal.GET("/HistoricalGoal", g.gc.HistoricalGoal)
		Goal.DELETE("/DeleteGoal/:task_id", g.gc.DeleteTask)
		Goal.POST("/CheckTask/:task_id", g.gc.CheckTask)
		Goal.POST("/MakeGoals", g.gc.PostGoals)
	}
}
