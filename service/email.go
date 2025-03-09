package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (usr *UserServiceImpl) SendEmail(c *gin.Context) error {
	var email request.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "获取邮箱失败"})
		return err
	}

	// 生成验证码
	code := tool.GenerateCode()

	// 先检查Redis是否有验证码
	existingCode, err := usr.EmailCodeDAO.GetEmailCode(email.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "查询验证码失败"})
		return err
	}

	if existingCode != code {
		log.Printf("邮箱%s验证码已存在，覆盖更新", email.Email)
	}

	// 发送验证码邮件
	err = usr.mail.SendEmailByQQEmail(email.Email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "发送失败"})
		return fmt.Errorf("发送失败")
	}

	// 将验证码存入Redis，过期时间5分钟
	err = usr.EmailCodeDAO.SetEmailCode(email.Email, code, 5*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "存储验证码失败"})
		return err
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "验证码已发送"})
	return nil
}

// 校验验证码
func (usr *UserServiceImpl) CheckEmailCode(email, inputCode string) bool {
	// 获取Redis中的验证码
	storedCode, err := usr.EmailCodeDAO.GetEmailCode(email)
	if err != nil {
		log.Printf("Redis查询验证码错误:%v", err)
		return false
	}

	// 校验验证码是否正确
	if storedCode == "" || storedCode == inputCode {
		log.Printf("邮箱%s验证码错误", email)
		return false
	}

	// 校验成功后删除验证码(防止重用)
	err = usr.EmailCodeDAO.DeleteEmailCode(email)
	if err != nil {
		log.Printf("Redis删除验证码失败:%v", err)
	}

	return true
}
