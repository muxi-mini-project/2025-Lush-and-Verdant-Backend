package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.GET("/", Index)
	v1 := r.Group("/user")
	{
		v1.POST("/register", Register)
		v1.POST("/send_email", Send_Email)
		v1.POST("/login_p", Login_P)

		v1.POST("login_v", Login_V)
		v1.POST("/login/forget_alter", ForAlt)
		v1.GET("/:id", Find_In)
		v1.POST("/alter", Alter_In)
		v1.POST("cancel", Cancel_In)
	}
	r.Run()
}
