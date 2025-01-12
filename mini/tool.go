package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)

const (
	pri       = "葱茏用户"
	secretKey = "aComplexSecretKeyForSecurity"
)

func Index(c *gin.Context) { //测试邮件功能
	code := GenerateCode()
	err := SendEmail("2085661244@qq.com", code)
	if err != nil {
		fmt.Println(err)
	}
}

// 定义一个结构体用于存储 Claims 信息
type Claims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(id int) (string, error) {
	claims := Claims{
		UserId: id, // 使用传入的用户 ID
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(), // 过期时间设置为 1 个月后
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey)) // 替换为你的密钥
}

// 创建一个中间件来验证 JWT 并检查用户角色：
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"message": "未认证"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil // 替换为你的密钥
		})

		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid token", "error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(401, gin.H{"message": "未认证"})
			c.Abort()
			return
		}

		// 获取 user_id 并转换为 int 类型
		userId, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(401, gin.H{"message": "无法获取id"})
			c.Abort()
			return
		}

		// 打印出用户ID，或存储在上下文中
		fmt.Println(userId)
		c.Set("user_id", int(userId)) // 将用户ID存入上下文中，方便后续使用

		c.Next()
	}
}

// @Summary 发送验证码邮件
// @Description 向用户邮箱发送验证码
// @Tags 邮件相关
// @Accept json
// @Produce json
// @Param email body Email true "邮箱发送验证码信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /send-email [post]
func Send_Email(c *gin.Context) {
	var email Email
	var count int = 0
	//读取前端发送的邮箱号
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	//生成验证码
	code := GenerateCode()
	email.Code = code
	//发送验证码
	err := SendEmail(email.Name, code)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = db.QueryRow("select count(*) from email where emailname = ?", email.Name).Scan(&count)
	if err != nil {
		c.JSON(400, Response{Error: err.Error()})

		return
	}
	if count == 0 {
		_, err = db.Exec("insert email (emailname, code, status) values (?, ?, ?)", email.Name, code, "valid")
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		_, err = db.Exec("update email set code = ?,status = ? where emailname = ?", code, "valid", email.Name)
		if err != nil {
			c.JSON(400, Response{Error: err.Error()})
			return
		}
	}

	//5分钟后删除验证码
	delay := 5 * time.Minute
	time.AfterFunc(delay, func() {
		_, err = db.Exec("update email set status=? where emailname = ? and status = ?", "expired", email.Name, "valid")
	})
	c.JSON(200, Response{Message: "发送成功"})

}

// @Summary 用户注册
// @Description 用户通过邮箱注册，包含验证码验证
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param user body UserRegister true "用户注册信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /register [post]
func Register(c *gin.Context) {
	var user UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, Response{Message: err.Error()})
	}
	var code_str string
	err := db.QueryRow("select code from email where emailname = ?", user.Email).Scan(&code_str)
	if err != nil {
		c.JSON(400, Response{Message: err.Error()})
	}

	code_num, err := strconv.Atoi(code_str)
	if err != nil {
		log.Println(err.Error())
	}
	userCode, err := strconv.Atoi(user.Code)
	if err != nil {
		log.Println(err.Error())
	}
	if code_num == userCode {
		log.Printf("用户 %s 验证成功！", user.Email)
	}
	_, err = db.Exec("insert  user (username,password,device_num,email) values (?, ?, ?, ?)", user.Username, user.Password, user.Device_Num, user.Email)
	if err != nil {
		c.JSON(400, Response{Message: err.Error()})
	}
}

// @Summary 用户登录（密码方式）
// @Description 用户通过邮箱和密码进行登录，成功后返回token
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param user body UserLogin true "用户登录信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /login/password [post]
func Login_P(c *gin.Context) { // 用户用密码登录
	var user UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, Response{Message: err.Error()})
	}

	var password string
	var id int
	err := db.QueryRow("select password,id from user where email = ?", user.Email).Scan(&password, &id)
	if err != nil {
		c.JSON(400, Response{Message: "未找到该用户" + err.Error()})
		log.Println(password, id, user.Email)
		return
	}

	//登录成功
	if user.Password == "" {
		c.JSON(400, Response{Message: "未输入密码"})
		return
	}

	if password == user.Password {
		//func GenerateToken(id int) (string, error) {
		token, err := GenerateToken(id)
		if err != err {
			c.JSON(401, Response{Message: err.Error()})
			return
		}
		c.JSON(200, Response{Message: "登录成功", Token: token})
	} else {
		c.JSON(400, Response{Message: "密码错误"})
		return
	}
}

// @Summary 游客登录
// @Description 游客通过设备号进行登录，如果该设备号尚未注册，则会创建一个新用户并返回 JWT Token
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param vister body VisterLogin true "游客登录信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /login/visitor [post]
func Login_V(c *gin.Context) { //游客登录
	var vister VisterLogin
	if err := c.ShouldBindJSON(&vister); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	_, err := db.Exec("insert user_v (device_num) values (?)", vister.Device_Num)
	//unique约束
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 { //重复
				c.JSON(200, gin.H{
					"msg": "用户已登录",
				})
				return
			} else {
				log.Println(err.Error())
				return
			}
		} else {
			log.Println(err.Error())
			return
		}
	} else {
		log.Printf("用户设备号为 %s 注册成功", vister.Device_Num)
	}
	var id int
	err = db.QueryRow("select id from user_v where device_num = ?", vister.Device_Num).Scan(&id)
	if err != nil {
		c.JSON(400, Response{Message: "登录出错" + err.Error()})
		return
	}
	//保存数据库
	username := pri + strconv.Itoa(id)
	_, err = db.Exec("update user_v set username = ? where device_num = ?", username, vister.Device_Num)
	token, err := GenerateToken(id)
	if err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	c.JSON(200, Response{Message: "游客登陆成功", Token: token})
}

// @Summary 忘记密码和修改密码
// @Description 用户通过邮箱和验证码修改密码，验证码有效后可以修改密码。
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param alter body AlterPassword true "修改密码信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /alter/password [post]
func ForAlt(c *gin.Context) { //忘记密码、修改密码
	var alter AlterPassword
	if err := c.ShouldBindJSON(&alter); err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}
	var code_str string
	err := db.QueryRow("select code from email where emailname = ?", alter.Email).Scan(&code_str)
	if err != nil {
		c.JSON(400, Response{Message: err.Error()})
		return
	}

	code_num, err := strconv.Atoi(code_str)
	if err != nil {
		log.Println(err.Error())
		return
	}

	code, err := strconv.Atoi(alter.Code)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if code_num == code { //验证成功
		_, err = db.Exec("update user set password = ? where email = ?", alter.Password, alter.Email)
		if err != nil {
			c.JSON(400, Response{Message: err.Error()})
			return
		} else {
			c.JSON(200, Response{Message: "修改成功"})
		}
	} else {
		c.JSON(400, Response{Message: "验证码错误"})
	}
}
func Find_In(c *gin.Context) {

}
func Alter_In(c *gin.Context) {

}

// @Summary 注销账户
// @Description 用户通过邮箱注销账户。如果邮箱对应的用户不存在，返回错误信息。
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param cancel body CancelUser true "注销账户信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /cancel/account [post]
func Cancel_In(c *gin.Context) { // 注销账户
	var cancel CancelUser
	if err := c.ShouldBindJSON(&cancel); err != nil {
		c.JSON(400, Response{Message: "请求参数无效: " + err.Error()})
		return
	}

	// 检查用户是否存在
	var count int
	err := db.QueryRow("select count(*) from user where email= ?", cancel.Email).Scan(&count)
	if err != nil {
		c.JSON(400, Response{Message: "数据库查询失败: " + err.Error()})
		return
	}
	if count == 0 {
		c.JSON(400, Response{Message: "用户不存在"})
		return
	}

	// 删除用户
	_, err = db.Exec("delete from user where email = ?", cancel.Email)
	if err != nil {
		c.JSON(400, Response{Message: "账户注销失败: " + err.Error()})
		return
	}

	// 记录注销操作日志
	log.Printf("用户 %s 注销成功", cancel.Email)

	// 返回成功消息
	c.JSON(200, Response{Message: "账户已成功注销"})
}
