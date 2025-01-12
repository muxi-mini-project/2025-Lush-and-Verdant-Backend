package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type Request struct {
	Newslogan string `json:"newslogan"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Token   string `json:"token,omitempty"`
}

type Slogans struct {
	Id     string `json:"id"`
	Slogan string `json:"slogan"`
}

type Users struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Device_Num string `json:"device_num"`
	Email      string `json:"email"`
	Goal_Pubic string `json:"goal_public"`
	Slogan     string `json:"slogan"`
}

type Newslogan struct {
	Newslogan string `json:"newslogan"`
	ID        string `json:"id"`
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
	return token.SignedString([]byte("1234567890")) // 替换为你的密钥
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

func Register(c *gin.Context) {
	var user Users
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "检查用户名出错"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, Response{Error: "用户已存在"})
		return
	}

	insert := "INSERT INTO users(username,password,role) VALUES(?, ?,?)"
	_, err = db.Exec(insert, user.Username, user.Password, "member")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: "注册失败"})
	}

	c.JSON(http.StatusCreated, Response{Message: "注册成功"})
}

func Login(c *gin.Context) {
	var user Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	var storedusers Users
	err := db.QueryRow("SELECT username,password,role FROM users WHERE username = ?", user.Username).Scan(&storedusers.Username, &storedusers.Password, &storedusers.Role)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Error: "用户不存在或密码错误"})
	}

	if storedusers.Password != user.Password {
		c.JSON(http.StatusBadRequest, Response{Error: "用户不存在或密码错误"})
		return
	}

	token, err := GenerateToken(storedusers.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "登录成功", Token: token})
}

func Change_words(c *gin.Context) {
	var newslogan Newslogan
	if err := c.ShouldBind(&newslogan); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, Response{Error: "无法获取用户ID"})
		return
	}

	id := userId.(int)

	_, err := db.Exec("UPDATE Users SET slogans = ? WHERE id = ?", newslogan.Newslogan, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "更新成功"})
}

func Get_words(c *gin.Context) {
	device := c.Param("device")

	if device == "" {
		c.JSON(http.StatusBadRequest, Response{Error: "设备号不能为空"})
		return
	}

	var exitslogan string
	err := db.QueryRow("SELECT slogan FROM Users WHERE device = ?", device).Scan(&exitslogan)
	if err == nil && exitslogan != "" {
		c.JSON(http.StatusAlreadyReported, Response{Message: "已拥有激励语"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(10) + 1

	var slogans Slogans
	err = db.QueryRow("SELECT id,slogan FROM slogan_list WHERE id = ?", id).Scan(&slogans.Id, &slogans.Slogan)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{Error: "读取失败"})
		return
	}

	_, err = db.Exec("UPDATE user SET slogan = ? WHERE device_num = ?", slogans.Slogan, device)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "更新失败"})
		return
	}

	c.JSON(http.StatusOK, Response{Message: "更新成功"})
}
