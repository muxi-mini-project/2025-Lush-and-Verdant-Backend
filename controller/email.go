package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (uc *UserController) SendEmail(c *gin.Context) {
	var email request.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: "获取邮箱失败"})
		return
	}

	err := uc.usr.SendEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "发送成功"})
}
