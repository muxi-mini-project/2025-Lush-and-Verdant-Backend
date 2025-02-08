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
		c.JSON(http.StatusBadRequest, response.Response{Error: "获取邮箱失败"})
		return err
	}

	//生成验证码
	code := tool.GenerateCode()
	//先查询验证码的状态
	emailCheck, ok := usr.CheckSendEmail(email.Email)
	if ok { //有验证码
		if emailCheck.Status { //此时是有效的，重新发送,并修改验证码
			err := usr.mail.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				return fmt.Errorf("发送失败")
			}
			//更新验证码
			result := usr.db.Model(&emailCheck).Where("mail = ?", email.Email).Update("code", code)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "更新验证码失败"})
				return fmt.Errorf("更新验证码失败")
			}
		} else { //此时是无效的，重新发送，并更新验证码及其状态
			err := usr.mail.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "发送失败"})
				return fmt.Errorf("发送失败")
			}
			result := usr.db.Model(&email).Where("mail = ?", email.Email).Updates(map[string]interface{}{"code": code, "status": true})
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "更新验证码失败"})
				return fmt.Errorf("更新验证码失败")
			}
		}
	} else { //没有验证码
		emailCheck.Status = true
		emailCheck.Code = code
		emailCheck.Email = email.Email
		err := usr.mail.SendEmailByQQEmail(email.Email, code)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: "发送失败"})
			return fmt.Errorf("发送失败")
		}
		result := usr.db.Create(&emailCheck)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: "更新验证码失败"})
			return fmt.Errorf("更新验证码失败")
		}
	}

	//因为不管怎么样都重新发送了，所以5分钟后都失效一波
	//5分钟后删除验证码
	delay := 5 * time.Minute
	time.AfterFunc(delay, func() {
		//先检查状态
		emailChe, _ := usr.CheckSendEmail(email.Email)
		//如果状态是有效的,变成无效的
		if emailChe.Status {
			result := usr.db.Model(&email).Where("mail = ?", email.Email).Update("status", false)
			if result.Error != nil {
				log.Println("用户%v的验证码状态改变出现错误", email.Email)
				return
			}
		}
		//状态无效则不做处理
	})
	c.JSON(http.StatusOK, response.Response{Message: "发送成功"})
	return nil

}
