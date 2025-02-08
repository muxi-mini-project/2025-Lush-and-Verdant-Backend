package service

import (
	"2025-Lush-and-Verdant-Backend/model"
	"github.com/gin-gonic/gin"
)

//创建接口

type UserService interface {
	UserRegister(*gin.Context) error
	CheckUserByDevice(string) (*model.User, bool)
	CheckSendEmail(string) (*model.Email, bool)
	UserLogin(*gin.Context) error
	CheckUserByEmail(string) (*model.User, bool)
	VisitorLogin(*gin.Context) error
	ForForAlt(*gin.Context) error
	Cancel(*gin.Context) error
	SendEmail(*gin.Context) error
}
