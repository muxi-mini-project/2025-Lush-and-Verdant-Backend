package service

import (
	"github.com/gin-gonic/gin"
)

//创建接口

type UserService interface {
	UserRegister(*gin.Context) error
	UserLogin(*gin.Context) error
	VisitorLogin(*gin.Context) error
	ForForAlt(*gin.Context) error
	Cancel(*gin.Context) error
	SendEmail(*gin.Context) error
}
