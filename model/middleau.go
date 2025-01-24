package model

import "github.com/dgrijalva/jwt-go"

// 定义一个结构体用于存储 Claims 信息
type Claims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}
