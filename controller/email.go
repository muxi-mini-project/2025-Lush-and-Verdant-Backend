package controller

import (
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (uc *UserController) SendEmail(c *gin.Context) {
	var email model.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: "获取邮箱失败"})
		return
	}
	//生成验证码
	code := tool.GenerateCode()

	db := dao.NewDB(dsn)

	//先查询验证码的状态
	result := db.Where("email = ?", email.Email).Find(&email)
	//如果有就检查验证码的状态
	if result.RowsAffected != 0 {
		if email.Status { //此时是有效的，重新发送,并修改验证码
			err := tool.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.Response{Error: "发送失败"})
				return
			}
			//更新验证码
			result = db.Model(&email).Where("email = ?", email.Email).Update("code", code)
			if result.Error != nil {
				c.JSON(http.StatusConflict, model.Response{Error: "更新验证码失败"})
				return
			}
		} else { //此时是无效的，重新发送，并更新验证码及其状态

			err := tool.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.Response{Error: "发送失败"})
				return
			}
			result := db.Model(&email).Where("email = ?", email.Email).Updates(map[string]interface{}{"code": code, "status": true})
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, model.Response{Error: "更新验证码失败"})
				return
			}
		}
	} else {
		email.Status = true
		email.Code = code
		err := tool.SendEmailByQQEmail(email.Email, code)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: "发送失败"})
			return
		}
		result := db.Create(&email)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.Response{Error: "更新验证码失败"})
			return
		}
	}
	//因为不管怎么样都重新发送了，所以5分钟后都失效一波
	//5分钟后删除验证码
	delay := 5 * time.Minute
	time.AfterFunc(delay, func() {
		//先检查状态
		result := db.Where("email = ?", email.Email).First(&email)
		if result.Error != nil {
			log.Println(result.Error.Error())
			return
		}
		//如果状态是有效的,变成无效的
		if email.Status {
			result := db.Model(&email).Where("email = ?", email.Email).Update("status", false)
			if result.Error != nil {
				log.Println(result.Error.Error())
				return
			}
		}
		//状态无效则不做处理
	})
	//发送成功
	c.JSON(200, model.Response{Message: "发送成功"})
}
