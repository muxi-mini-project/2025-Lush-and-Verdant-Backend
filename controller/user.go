package controller

import (
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"log"
)

type UserController struct {
	usr service.UserService
}

// 创建新用户
func NewUserController(usr service.UserService) *UserController {
	return &UserController{
		usr: usr,
	}
}

// 注册用户
func (uc *UserController) Register(c *gin.Context) {
	err := uc.usr.UserRegister(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// 登录用户(邮箱和密码登录)(正式用户)
func (uc *UserController) Login(c *gin.Context) {
	err := uc.usr.UserLogin(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// 游客登录
func (uc *UserController) Login_v(c *gin.Context) {
	err := uc.usr.VisitorLogin(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// 忘记密码和修改密码
func (uc *UserController) ForAlt(c *gin.Context) {
	err := uc.usr.ForForAlt(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// 用户注销
func (uc *UserController) Cancel(c *gin.Context) {
	err := uc.usr.Cancel(c)
	if err != nil {
		log.Println(err)
		return
	}
}
