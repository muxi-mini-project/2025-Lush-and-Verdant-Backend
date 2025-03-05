package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

// SendEmail 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 通过用户提供的邮箱发送验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param email body string true "用户邮箱"
// @Success 200 {string} string "验证码发送成功"
// @Failure 400 {string} string "请求错误"
// @Failure 500 {string} string "服务器错误"
// @Router /user/send_email [post]
func (uc *UserController) SendEmail(c *gin.Context) {
	err := uc.usr.SendEmail(c)
	if err != nil {
		log.Println(err)
		return
	}
}
