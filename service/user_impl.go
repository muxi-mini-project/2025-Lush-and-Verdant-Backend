package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/api/response"
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/middleware"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserServiceImpl struct {
	Dao    dao.UserDAO
	jwt    *middleware.JwtClient
	mail   *tool.Mail
	priCfg *config.PriConfig
}

func NewUserServiceImpl(userDao dao.UserDAO, jwt *middleware.JwtClient, mail *tool.Mail, priCfg *config.PriConfig) *UserServiceImpl {
	return &UserServiceImpl{
		Dao:    userDao,
		jwt:    jwt,
		mail:   mail,
		priCfg: priCfg,
	}
}

func (usr *UserServiceImpl) UserRegister(c *gin.Context) error {
	//获取前端的消息
	var userRegister request.UserRegister
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return err
	}
	//先查询用户是否注册
	user, ok := usr.Dao.CheckUserByDevice(userRegister.Device_Num)
	if ok { //有用户
		if strings.Compare(user.Password, "") == 0 { //此时是游客，更新状态
			//此时更新用户
			user.Username = userRegister.Username
			user.Password = userRegister.Password
			user.Email = userRegister.Email
			//更新用户
			err := usr.Dao.VisitorToUser(userRegister.Device_Num, user)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "游客转正失败"})
				return err
			}

			log.Printf("register user %v success", user.Email)
			err = usr.mail.ChangeStatus(userRegister.Email)
			if err != nil {
				log.Println("转正完成但验证码未修改成功")
				log.Println(err)
				c.JSON(http.StatusBadRequest, response.Response{Error: "转正完成但验证码未修改成功"})
				return err
			}
			c.JSON(http.StatusOK, response.Response{Message: "游客转正成功"})
			return nil
		} else {
			c.JSON(http.StatusTooManyRequests, response.Response{Error: "用户已注册"})
			return fmt.Errorf("用户已注册")
		}
	} else { //没有用户
		//先查询验证码的状态
		email, ok := usr.Dao.CheckSendEmail(userRegister.Email)
		if ok { //有验证码
			if email.Status { //验证码有效
				if strings.Compare(email.Code, userRegister.Code) == 0 { //验证码验证成功
					//此时注册用户
					user.Username = userRegister.Username
					user.Password = userRegister.Password
					user.Email = userRegister.Email
					user.DeviceNum = userRegister.Device_Num
					//创建新用户
					err := usr.Dao.CreateUser(user)
					if err != nil {
						c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
						return err
					}
					log.Printf("register user %v success", user.Email)
					err = usr.mail.ChangeStatus(userRegister.Email)
					if err != nil {
						log.Println("注册完成但验证码未修改成功")
						log.Println(err)
						c.JSON(http.StatusBadRequest, response.Response{Error: "注册完成但验证码未修改成功"})
						return err
					}
					c.JSON(http.StatusOK, response.Response{Message: "用户注册成功"})
					return nil //用户注册成功

				} else {
					c.JSON(http.StatusBadRequest, response.Response{Error: "验证码错误"})
					return fmt.Errorf("验证码错误")
				}
			} else {
				c.JSON(http.StatusOK, response.Response{Error: "验证码无效"})
				return fmt.Errorf("验证码无效")
			}
		} else {
			c.JSON(http.StatusNotFound, response.Response{Error: "用户未发送验证码"})
			return fmt.Errorf("用户未发送验证码")
		}
	}
}

// 用户登录
func (usr *UserServiceImpl) UserLogin(c *gin.Context) error {
	//接受前端发送的消息
	var userLogin request.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return err
	}

	user, ok := usr.Dao.CheckUserByEmail(userLogin.Email)
	if ok { //发现用户，验证密码
		if user.Password == userLogin.Password { //密码正确
			//生成token
			token, err := usr.jwt.GenerateToken(int(user.ID))
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: "生成token失败"})
				return err
			}
			c.JSON(http.StatusOK, response.Response{Message: "登录成功", Token: token})
			return nil
		} else {
			c.JSON(http.StatusConflict, response.Response{Error: "密码错误"})
			return fmt.Errorf("%s 密码错误", user.Email)
		}
	} else { //没有该用户
		c.JSON(http.StatusBadRequest, response.Response{Error: "用户未注册"})
		return fmt.Errorf("%s 用户未注册", userLogin.Email)
	}

}

// 检查游客注册和登录
func (usr *UserServiceImpl) VisitorLogin(c *gin.Context) error {
	//接受前端的消息
	var visitor request.Visitor
	if err := c.ShouldBindJSON(&visitor); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return err
	}
	//查询这个机型是否登陆过
	user, ok := usr.Dao.CheckUserByDevice(visitor.Device_Num)
	if ok { //登录过 //不是新用户
		now := time.Now()
		monthlate := user.CreatedAt.AddDate(0, 1, 0)
		timeSub := monthlate.Sub(now)
		// 输出时间差（天、小时、分钟、秒）
		days := int(timeSub.Hours()) / 24
		hours := int(timeSub.Hours()) % 24
		minutes := int(timeSub.Minutes()) % 60
		seconds := int(timeSub.Seconds()) % 60
		duration := fmt.Sprintf("%02d天-%02d时:%02d分:%02d秒", days, hours, minutes, seconds)
		if now.After(monthlate) { //超过一个月了
			c.JSON(http.StatusConflict, response.Response{Error: "游客登录时间过长,禁止登录"})
			return fmt.Errorf("%s 游客登录时间过长,禁止登录", user.DeviceNum)
		} else { //没超过一个月
			token, err := usr.jwt.GenerateToken(int(user.ID))
			if err != nil {
				c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
				return err
			}
			c.JSON(http.StatusOK, response.Response{Message: "游客登录成功,剩余时间为" + duration, Token: token})
			return nil
		}
	} else { //说明是新用户
		user.DeviceNum = visitor.Device_Num
		//创建游客
		err := usr.Dao.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
			return err
		}
		//给这个游客姓名
		username := usr.priCfg.Name + strconv.Itoa(int(user.ID))
		err = usr.Dao.UpdateUserName(user.DeviceNum, username)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
			return err
		}
		//传送token
		token, err := usr.jwt.GenerateToken(int(user.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
			return err
		}
		c.JSON(http.StatusOK, response.Response{Message: "游客注册成功，剩余时间30天", Token: token})
		return nil
	}
}

// 忘记密码 -> 修改密码
func (usr *UserServiceImpl) ForForAlt(c *gin.Context) error {
	//接受前端消息
	var foralt request.ForAlter
	if err := c.ShouldBindJSON(&foralt); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return err
	}

	//先查有没有该用户
	user, ok := usr.Dao.CheckUserByEmail(foralt.Email)
	if ok { //有用户
		//查询验证码
		email, ok := usr.Dao.CheckSendEmail(foralt.Email)
		if ok { //查询到验证码
			if email.Status { //验证码有效
				if strings.Compare(email.Code, foralt.Code) == 0 { //验证码正确
					err := usr.Dao.UpdatePassword(foralt.Email, foralt.Password)
					if err != nil {
						c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
						return err
					}
					err = usr.mail.ChangeStatus(email.Email)
					if err != nil {
						c.JSON(http.StatusBadRequest, response.Response{Error: "修改密码成功但是验证码状态设置失败"})
						return fmt.Errorf("%s 修改密码成功但是验证码状态设置失败:%s", user.Email, err.Error())
					}
					c.JSON(http.StatusOK, response.Response{Message: "修改密码成功"})
					return nil
				} else {
					c.JSON(http.StatusConflict, response.Response{Error: "验证码错误"})
					return fmt.Errorf("验证码错误")
				}
			} else { //验证码失效
				c.JSON(http.StatusBadRequest, response.Response{Error: "验证码失效"})
				return fmt.Errorf("验证码失效")
			}
		} else { //未查询到验证码
			c.JSON(http.StatusBadRequest, response.Response{Error: "没有发验证码"})
			return fmt.Errorf("没有发验证码")
		}
	} else { //没有用户
		c.JSON(http.StatusNotFound, response.Response{Error: "没有该用户"})
		return fmt.Errorf("没有该用户")

	}

}

// 用户注销
func (usr *UserServiceImpl) Cancel(c *gin.Context) error {
	//接受前端消息
	var cancel request.UserCancel
	if err := c.ShouldBindJSON(&cancel); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
		return err
	}
	//查询用户
	user, ok := usr.Dao.CheckUserByEmail(cancel.Email)
	if ok { //找到了,直接硬删除
		err := usr.Dao.DeleteUser(cancel.Email, user)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
			return err
		}
		c.JSON(http.StatusOK, response.Response{Message: "用户注销成功"})
		return nil
	} else {
		c.JSON(http.StatusNotFound, response.Response{Error: "用户没注册，还妄想注销账户"})
		return fmt.Errorf("%s 用户没注册，还妄想注销账户", cancel.Email)
	}
}
