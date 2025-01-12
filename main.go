package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Courage_words struct {
	words string
}

func main() {
	InitDB()
	if err := InitDB(); err != nil {
		log.Fatal("无法连接到数据库")
	}

	r := gin.Default()

	r.PUT("/courage_words/change_word", AuthMiddleware(), Change_words)
	r.GET("/courage_words/get_word/:device", Get_words)

	r.Run(":8080")
}
