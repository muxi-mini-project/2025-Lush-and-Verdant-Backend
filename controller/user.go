package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var dsn = config.Dsn

type UserController struct {
	usr *service.UserService
}

// 创建新用户
func NewUserController(usr *service.UserService) *UserController {
	return &UserController{
		usr: usr,
	}
}

// 注册用户
func (uc *UserController) Register(c *gin.Context) {
	//获取前端的消息
	var userRegister request.UserRegister
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	code, err := uc.usr.CheckUserRegister(userRegister)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	switch code {
	case 1: //用户注册成功
		c.JSON(http.StatusOK, response.Response{Message: "用户注册成功"})
	case 2: //游客转正成功
		c.JSON(http.StatusOK, response.Response{Message: "游客转正成功"})
	}
}

// 登录用户(邮箱和密码登录)(正式用户)
func (uc *UserController) Login(c *gin.Context) {
	var userLogin request.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	token, err := uc.usr.CheckUserLogin(userLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Message: "登陆成功", Token: token})

}

// 游客登录
func (uc *UserController) Login_v(c *gin.Context) {
	var visiter request.Visiter
	if err := c.ShouldBindJSON(&visiter); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	token, msg, err := uc.usr.CheckVisiterLogin(visiter)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: msg, Token: token})
}

// 忘记密码和修改密码
func (uc *UserController) ForAlt(c *gin.Context) {
	var foralt request.ForAlter
	if err := c.ShouldBindJSON(&foralt); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	err := uc.usr.ForForAlt(foralt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	}
	c.JSON(http.StatusOK, response.Response{Message: "修改密码成功"})
}

// 用户注销
func (uc *UserController) Cancel(c *gin.Context) {
	var cancel request.UserCancel
	if err := c.ShouldBindJSON(&cancel); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	err := uc.usr.Cancel(cancel)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Message: "注销成功"})
}
