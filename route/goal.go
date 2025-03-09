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

	Goal := r.Group("/goal").Use(g.jwt.AuthMiddleware())
	{
		Goal.GET("/GetGoal", g.gc.GetGoal)
		Goal.POST("/MakeGoal", g.gc.PostGoal)
		Goal.PUT("/UpdateGoal/:goal_id", g.gc.UpdateGoal)
		Goal.GET("/HistoricalGoal", g.gc.HistoricalGoal)
		Goal.DELETE("/DeleteGoal/:goal_id", g.gc.DeleteGoal)
	}
}
