package tool

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/smtp"
	"time"
)

var dsn = config.Dsn

// 生成验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano()) //以纳米为级别
	code := rand.Intn(1000000)       //生成6位数的验证码
	return fmt.Sprintf("%06d", code)
}

// SendEmailByQQEmail 发送邮件函数
func SendEmailByQQEmail(to string, code string) error {
	from := "3953017473@qq.com"
	password := "vzsvxefmdmqkcgbg" // 邮箱授权码
	smtpServer := "smtp.qq.com:465"

	// 设置 PlainAuth
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")

	// 创建 tls 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()

	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()

	subject := "Lush-And-Verdant"
	body := `
		<h1>Verification Code</h1>
		<p>Your verification code is: <strong>` + code + `</strong></p >
		<p>This verification code is valid for 5 minutes</p >
		<p>If you are not doing it yourself, please ignore it !</p >
		<h1>验证码</h1>
		<p>你的验证码是: <strong>` + code + `</strong></p >
		<p>验证码的有效时间是5分钟。</p >
		<p>如非本人操作，请忽略此邮件！</p >
	`
	msg := []byte("From: Sender Name <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body)

	_, err = wc.Write(msg)
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}

	return nil
}

// 从前端获得email json格式的
func GetEmailName(c *gin.Context) (string, bool) {
	var email model.Email
	if err := c.ShouldBind(&email); err != nil {
		return "", false
	} else {
		return email.Email, true
	}
}

// 修改验证码状态
func ChangeStatus(email string) error {
	db := dao.NewDB(dsn)
	var email1 model.Email
	//查询信息

	result := db.Where("email=?", email).First(&email1)
	if result.Error != nil {
		return result.Error
	}
	if email1.Status { //有效改为无效
		result := db.Where("email=?", email).Update("status", false)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		return nil
	}

}
