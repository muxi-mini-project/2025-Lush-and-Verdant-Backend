package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"log"
	"time"
)

func (usr *UserService) SendEmail(email request.Email) error {

	//生成验证码
	code := tool.GenerateCode()

	var emailCheck model.Email
	//先查询验证码的状态
	result := usr.db.Where("email = ?", email.Email).Find(&emailCheck)
	//如果有就检查验证码的状态
	if result.RowsAffected != 0 {
		if emailCheck.Status { //此时是有效的，重新发送,并修改验证码
			err := tool.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				return fmt.Errorf("发送失败")
			}
			//更新验证码
			result = usr.db.Model(&emailCheck).Where("email = ?", email.Email).Update("code", code)
			if result.Error != nil {
				return fmt.Errorf("更新验证码失败")
			}
		} else { //此时是无效的，重新发送，并更新验证码及其状态

			err := tool.SendEmailByQQEmail(email.Email, code)
			if err != nil {
				return fmt.Errorf("发送失败")
			}
			result := usr.db.Model(&email).Where("email = ?", email.Email).Updates(map[string]interface{}{"code": code, "status": true})
			if result.Error != nil {
				return fmt.Errorf("更新验证码失败")
			}
		}
	} else { //没有
		emailCheck.Status = true
		emailCheck.Code = code
		emailCheck.Email = email.Email
		err := tool.SendEmailByQQEmail(email.Email, code)
		if err != nil {
			return fmt.Errorf("发送失败")
		}
		result := usr.db.Create(&emailCheck)
		if result.Error != nil {
			return fmt.Errorf("更新验证码失败")
		}
	}

	//因为不管怎么样都重新发送了，所以5分钟后都失效一波
	//5分钟后删除验证码
	delay := 5 * time.Minute
	time.AfterFunc(delay, func() {
		//先检查状态
		result := usr.db.Where("email = ?", email.Email).First(&emailCheck)
		if result.Error != nil {
			log.Println("用户%v的验证码查询出现错误", email.Email)
			return
		}
		//如果状态是有效的,变成无效的
		if emailCheck.Status {
			result := usr.db.Model(&email).Where("email = ?", email.Email).Update("status", false)
			if result.Error != nil {
				log.Println("用户%v的验证码状态改变出现错误", email.Email)
				return
			}
		}
		//状态无效则不做处理
	})

	return nil

}
