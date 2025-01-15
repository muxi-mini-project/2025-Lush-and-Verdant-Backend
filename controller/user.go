package controller

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/middleware"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

var dsn = config.Dsn

type UserController struct {
}

// 创建新用户
func NewUserController() *UserController {
	return &UserController{}
}

// 注册用户
func (uc *UserController) Register(c *gin.Context) {
	var user model.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})

		return
	}
	db := dao.NewDB(dsn)
	//先查询是否有用户
	var usercheck model.User
	result := db.Where("email = ?", user.Email).Find(&usercheck)
	if result.RowsAffected == 0 { //没有此用户
		//先查询验证码
		var email model.Email

		result := db.Where("email = ?", user.Email).First(&email)

		fmt.Println(";;;;;", email)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})

			return
		}
		if email.Status { //验证码有效
			if strings.Compare(email.Code, user.Code) == 0 { //验证码验证成功
				//此时注册用户
				usercheck.Username = user.Username
				usercheck.Password = user.Password
				usercheck.Email = user.Email
				usercheck.DeviceNum = user.Device_Num
				result := db.Create(&usercheck)
				if result.Error != nil {
					c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})

					return
				}
				log.Printf("register user %v success", usercheck.Email)
				c.JSON(http.StatusOK, model.Response{Message: "注册成功"})
				err := tool.ChangeStatus(user.Email)
				if err != nil {
					log.Println("注册完成但验证码未修改成功")
					log.Println(err)
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, model.Response{Message: "验证码失效"})
				return
			}
		}
	} else { //已经有用户了
		c.JSON(http.StatusConflict, model.Response{Message: "用户已存在"})
		return
	}

}

// 登录用户(邮箱和密码登录)(正式用户)
func (uc *UserController) Login(c *gin.Context) {
	var userlogin model.UserLogin
	if err := c.ShouldBindJSON(&userlogin); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
		return
	}
	db := dao.NewDB(dsn)
	var user model.User
	result := db.Where("email = ?", userlogin.Email).Find(&user)
	if result.RowsAffected == 0 { //没有该用户
		c.JSON(http.StatusBadRequest, model.Response{Message: "用户未注册"})
		return
	} else { //发现用户，验证密码
		if user.Password == userlogin.Password { //密码正确
			//生成token
			token, err := middleware.GenerateToken(int(user.ID))
			if err != nil {
				c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
				return
			}
			c.JSON(http.StatusOK, model.Response{Message: "登录成功", Token: token})
		} else {
			c.JSON(http.StatusBadRequest, model.Response{Message: "密码错误"})
			return
		}
	}
}

// 游客登录
func (uc *UserController) Login_v(c *gin.Context) {
	var visiter model.Visiter
	if err := c.ShouldBindJSON(&visiter); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
		return
	}
	db := dao.NewDB(dsn)
	var user model.User
	//查询这个机型是否登陆过
	result := db.Where("device_num = ?", visiter.Device_Num).Find(&user)
	if result.RowsAffected == 0 { //说明这是个新用户
		user.DeviceNum = visiter.Device_Num
		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})
			return
		}
		token, err := middleware.GenerateToken(int(user.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.Response{Message: "游客登录成功", Token: token})
	} else {
		now := time.Now()

		monthlate := user.CreatedAt.AddDate(0, 1, 0)
		if now.After(monthlate) { //超过一个月了
			c.JSON(http.StatusConflict, model.Response{Message: "游客登录时间过长"})
			return
		}
	}
}

// 忘记密码和修改密码
func (uc *UserController) ForAlt(c *gin.Context) {
	var foralt model.ForAlter
	if err := c.ShouldBindJSON(&foralt); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
		return
	}
	db := dao.NewDB(dsn)

	var user model.User

	//先查有没有该用户
	result := db.Where("email = ?", foralt.Email).Find(&user)
	if result.RowsAffected == 0 { //没有该用户
		c.JSON(http.StatusNotFound, model.Response{Message: "该用户未注册"})
		return
	} else {
		//查询验证码
		var email model.Email
		result := db.Where("email = ?", foralt.Email).Find(&email)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})
			return
		}
		if email.Status { //验证码有效
			if strings.Compare(email.Code, foralt.Code) == 0 { //验证码正确
				var user model.User
				result := db.Model(&user).Where("email = ?", foralt.Email).Update("password", foralt.Password)
				if result.Error != nil {
					c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})
					return
				}
				c.JSON(http.StatusOK, model.Response{Message: "修改成功"})
				err := tool.ChangeStatus(email.Email)
				if err != nil {
					log.Println("修改密码成功但是验证码状态设置失败")
					log.Println(err)
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, model.Response{Message: "验证码错误"})
				return
			}
		} else { //验证码失效
			c.JSON(http.StatusConflict, model.Response{Message: "验证码失效"})
			return
		}
	}
}

// 用户注销
func (uc *UserController) Cancel(c *gin.Context) {
	var cancel model.UserCancel
	if err := c.ShouldBindJSON(&cancel); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Error: err.Error()})
		return
	}
	db := dao.NewDB(dsn)
	var user model.User
	result := db.Where("email = ?", cancel.Email).Find(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, model.Response{Message: "用户没注册，还妄想注销账户"})
		return
	} else { //找到了,直接硬删除
		result = db.Model(&user).Where("email = ?", cancel.Email).Delete(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.Response{Error: result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, model.Response{Message: "用户注销成功"})
	}

}
