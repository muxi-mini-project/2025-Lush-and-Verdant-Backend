package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/middleware"
	"2025-Lush-and-Verdant-Backend/model"
	"2025-Lush-and-Verdant-Backend/tool"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// 此时的0, 1, 2相当于状态码
// 0 失败
// 1 注册
// 2 游客转正
func (usr *UserService) CheckUserRegister(userRegister request.UserRegister) (int, error) {
	//先查询是否有用户
	var user model.User
	result := usr.db.Where("device_num = ?", userRegister.Device_Num).First(&user)
	if result.RowsAffected == 0 { //没有此用户
		//先查询验证码的状态
		var email model.Email
		result := usr.db.Where("email = ?", userRegister.Email).First(&email)
		if result.RowsAffected == 0 { //用户未发送验证码
			return 0, fmt.Errorf("用户未发送验证码")
		} else {
			if email.Status { //验证码有效
				if strings.Compare(email.Code, userRegister.Code) == 0 { //验证码验证成功
					//此时注册用户
					user.Username = userRegister.Username
					user.Password = userRegister.Password
					user.Email = userRegister.Email
					user.DeviceNum = userRegister.Device_Num
					result := usr.db.Create(&user)
					if result.Error != nil {
						return 0, result.Error
					}
					log.Printf("register user %v success", user.Email)
					err := tool.ChangeStatus(userRegister.Email)
					if err != nil {
						log.Println("注册完成但验证码未修改成功")
						log.Println(err)
						return 0, err
					}
					return 1, nil //用户注册成功
				} else {
					return 0, fmt.Errorf("验证码错误")
				}
			} else {
				return 0, fmt.Errorf("验证码无效")
			}

		}

	} else { //已经有用户了，查看是游客还是用户
		if strings.Compare(user.Password, "") == 0 { //此时是游客，更新状态
			//此时更新用户
			user.Username = userRegister.Username
			user.Password = userRegister.Password
			user.Email = userRegister.Email
			result := usr.db.Model(&user).Where("device_num = ?", userRegister.Device_Num).Updates(&user)
			if result.Error != nil {
				return 0, result.Error
			}
			log.Printf("register user %v success", user.Email)
			err := tool.ChangeStatus(userRegister.Email)
			if err != nil {
				log.Println("转正完成但验证码未修改成功")
				log.Println(err)
				return 0, err
			}
			return 2, nil
		} else {
			return 0, fmt.Errorf("用户已注册")
		}
	}

}

// 检查用户登录
func (usr *UserService) CheckUserLogin(userLogin request.UserLogin) (string, error) {
	var user model.User
	result := usr.db.Where("email = ?", userLogin.Email).Find(&user)
	if result.RowsAffected == 0 { //没有该用户
		return "", fmt.Errorf("用户未注册")
	} else { //发现用户，验证密码
		if user.Password == userLogin.Password { //密码正确
			//生成token
			token, err := middleware.GenerateToken(int(user.ID))
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			return "", fmt.Errorf("密码错误")
		}
	}
}

// 检查游客注册和登录
func (usr *UserService) CheckVisiterLogin(visiter request.Visiter) (string, string, error) {

	var user model.User
	//查询这个机型是否登陆过
	result := usr.db.Model(&user).Where("device_num = ?", visiter.Device_Num).First(&user)
	if result.RowsAffected == 0 { //说明这是个新用户
		user.DeviceNum = visiter.Device_Num
		result := usr.db.Create(&user)
		if result.Error != nil {
			return "", "", result.Error
		}
		//给这个游客姓名
		username := config.Pri + strconv.Itoa(int(user.ID))
		result = usr.db.Model(&user).Where("device_num = ?", user.DeviceNum).Update("username", username)
		if result.Error != nil {
			return "", "", result.Error
		}
		//传送token
		token, err := middleware.GenerateToken(int(user.ID))
		if err != nil {
			return "", "", err
		}
		return token, "游客注册成功，剩余时间30天", nil
	} else { //不是新用户
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
			return "", "", fmt.Errorf("游客登录时间过长,禁止登录")
		} else { //没超过一个月
			token, err := middleware.GenerateToken(int(user.ID))
			if err != nil {
				return "", "", err
			}
			return token, "游客登录成功,剩余时间为" + duration, nil
		}

	}
}

func (usr *UserService) ForForAlt(foralt request.ForAlter) error {

	var user model.User
	//先查有没有该用户
	result := usr.db.Where("email = ?", foralt.Email).Find(&user)
	if result.RowsAffected == 0 { //没有该用户
		return fmt.Errorf("没有该用户")
	} else {
		//查询验证码
		var email model.Email
		result := usr.db.Where("email = ?", foralt.Email).Find(&email)
		if result.RowsAffected == 0 { //没有发验证码
			return fmt.Errorf("没有发验证码")
		} else {
			if email.Status { //验证码有效
				if strings.Compare(email.Code, foralt.Code) == 0 { //验证码正确
					var user model.User
					result := usr.db.Model(&user).Where("email = ?", foralt.Email).Update("password", foralt.Password)
					if result.Error != nil {
						return result.Error
					}
					err := tool.ChangeStatus(email.Email)
					if err != nil {
						log.Println("修改密码成功但是验证码状态设置失败")
						log.Println(err)
						return err
					}
					return nil
				} else {
					return fmt.Errorf("验证码错误")
				}
			} else { //验证码失效
				return fmt.Errorf("验证码失效")
			}
		}
	}
}

func (usr *UserService) Cancel(cancel request.UserCancel) error {

	var user model.User
	result := usr.db.Where("email = ?", cancel.Email).Find(&user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户没注册，还妄想注销账户")
	} else { //找到了,直接硬删除
		result = usr.db.Model(&user).Where("email = ?", cancel.Email).Unscoped().Delete(&user)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}
