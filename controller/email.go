package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (uc *UserController) SendEmail(c *gin.Context) {
	err := uc.usr.SendEmail(c)
	if err != nil {
		log.Println(err)
		return
	}
}
