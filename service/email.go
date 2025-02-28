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
	emailCheck, ok := usr.Dao.CheckSendEmail(email.Email)
	if ok { //有验证码
		if emailCheck.Status { //此时是有效的，重新发送,并修改验证码
			err := usr.mail.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				return fmt.Errorf("发送失败")
			}
			//更新验证码
			emailCheck.Code = code
			err = usr.Dao.UpdateUserEmail(emailCheck)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
				return err
			}

		} else { //此时是无效的，重新发送，并更新验证码及其状态
			err := usr.mail.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "发送失败"})
				return fmt.Errorf("发送失败")
			}

			emailCheck.Code = code
			emailCheck.Status = true
			err = usr.Dao.UpdateUserEmail(emailCheck)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
				return err
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

		err = usr.Dao.CreateUserEmail(emailCheck)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
			return err
		}
	}

	//因为不管怎么样都重新发送了，所以5分钟后都失效一波
	//5分钟后删除验证码
	delay := 5 * time.Minute
	time.AfterFunc(delay, func() {
		//先检查状态
		emailChe, _ := usr.Dao.CheckSendEmail(email.Email)
		//如果状态是有效的,变成无效的
		if emailChe.Status {
			emailChe.Status = false
			err := usr.Dao.UpdateUserEmail(emailChe)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
				log.Printf("用户%v的验证码状态改变出现错误\n", emailChe.Email)
				return
			}
		}
		//状态无效则不做处理
	})
	c.JSON(http.StatusOK, response.Response{Message: "发送成功"})
	return nil

}
