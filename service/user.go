package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/model"
	"github.com/gin-gonic/gin"
)

// 创建接口
type UserService interface {
	UserRegister(*gin.Context) error
	UserLogin(*gin.Context) error
	VisitorLogin(*gin.Context) error
	ForForAlt(*gin.Context) error
	Cancel(*gin.Context) error
	SendEmail(*gin.Context) error
	GetUserInfoById(idStr string) (*response.User, error)
	UpdateUserInfo(user *request.UserUpdate) error
	RandUser() (*model.User, error)
}
