package main

import "database/sql"

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Device_Num string `json:"device_num"`
	Email      string `json:"email"`
	Goal_Pubic string `json:"goal_public"`
	Slogans    string `json:"slogans"`
}
type Email struct {
	EmailName string `json:"email_name"`
	Code      string `json:"code"`
}
type UserRegister struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Device_Num string `json:"device_num"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}
type VisiterLogin struct {
	Username   string `json:"username"`
	Device_Num string `json:"device_num"`
	Email      string `json:"email"`
}
type AlterPassword struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
type CancelUser struct {
	Email string `json:"email"`
}

var db *sql.DB
