package controller

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

// Register 用户注册
// @Summary 用户注册
// @Description 通过邮箱和密码注册新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body request.UserRegister true "用户注册信息"
// @Success 200 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	err := uc.usr.UserRegister(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// Login 用户登录（正式用户）
// @Summary 用户登录
// @Description 通过邮箱和密码登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param credentials body request.UserLogin true "用户登录信息"
// @Success 200 {object} response.Response{data=response.Token} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	err := uc.usr.UserLogin(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// Login_v 游客登录
// @Summary 游客登录
// @Description 允许游客模式登录，无需账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body request.Visitor true "用户游客登陆"
// @Success 200 {object} response.Response{data=response.Token} "游客登录成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/login_v [post]
func (uc *UserController) Login_v(c *gin.Context) {
	err := uc.usr.VisitorLogin(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// ForAlt 忘记密码和修改密码
// @Summary 忘记密码/修改密码
// @Description 允许用户找回或更改密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param resetData body request.ForAlter true "密码修改信息"
// @Success 200 {object} response.Response "密码修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "验证码错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/foralt [post]
func (uc *UserController) ForAlt(c *gin.Context) {
	err := uc.usr.ForForAlt(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// Cancel 用户注销
// @Summary 用户注销
// @Description 允许用户注销账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body request.UserCancel true "注销信息 {email: string, password: string}"
// @Success 200 {object} response.Response "注销成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/cancel [post]
func (uc *UserController) Cancel(c *gin.Context) {
	err := uc.usr.Cancel(c)
	if err != nil {
		log.Println(err)
		return
	}
}

// GetUserInfoById 获取用户个人信息
// @Summary 获取用户个人信息
// @Description 允许用户获取个人信息
// @Tags 用户
// @Produce json
// @Success 200 {object} response.Response{data=response.User} "获取个人信息成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/info/{id} [GET]
func (uc *UserController) GetUserInfoById(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.usr.GetUserInfoById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "获取失败"})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "获取用户信息成功", Data: user})
}

// UpdateUserInfo 更新用户个人信息（姓名和邮箱）
// @Summary 更新用户个人信息
// @Description 允许用户更新个人信息
// @Tags 用户
// @Produce json
// @Success 200 {object} response.Response "更新个人信息成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /user/update [POST]
func (uc *UserController) UpdateUserInfo(c *gin.Context) {
	var userUpdate request.UserUpdate
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "参数传输失败"})
		return
	}
	err := uc.usr.UpdateUserInfo(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "修改成功"})
}
