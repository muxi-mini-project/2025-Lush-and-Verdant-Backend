package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type UserDAO interface {
	GetUserById(id uint) (*model.User, error)
	CheckUserByEmail(string) (*model.User, bool)
	CheckSendEmail(string) (*model.Email, bool)
	CheckUserByDevice(string) (*model.User, bool)
	VisitorToUser(string, *model.User) error
	CreateUser(*model.User) error
	CreateUserEmail(email *model.Email) error
	UpdateUser(user *model.User) error
	UpdateUserName(device string, username string) error
	UpdateUserNameById(id uint, name string) error
	UpdateUserEmailById(id uint, email string) error
	UpdatePassword(email string, password string) error
	UpdateUserEmail(email *model.Email) error
	DeleteUser(email string, user *model.User) error
	RandUser() (*model.User, error)
}

type UserDAOImpl struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAOImpl {
	return &UserDAOImpl{db: db}
}

// 通过用户id来查找用户
func (dao *UserDAOImpl) GetUserById(id uint) (*model.User, error) {
	var user model.User
	result := dao.db.Table("users").Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 通过邮箱查找是否有用户，并返回用户
func (dao *UserDAOImpl) CheckUserByEmail(addr string) (*model.User, bool) {
	var user model.User
	result := dao.db.Where("email = ?", addr).Find(&user)
	if result.RowsAffected == 0 {
		return &user, false
	} else {
		return &user, true
	}
}

// 检查是否发送验证码，并返回email
func (dao *UserDAOImpl) CheckSendEmail(addr string) (*model.Email, bool) {
	var email model.Email
	result := dao.db.Where("email = ?", addr).First(&email)
	if result.RowsAffected == 0 { //用户未发送验证码
		return &email, false
	} else {
		return &email, true
	}
}

// 通过设备号，检查是否拥有用户，并返回用户
func (dao *UserDAOImpl) CheckUserByDevice(deviceNum string) (*model.User, bool) {
	var user model.User
	result := dao.db.Where("device_num = ?", deviceNum).First(&user)
	if result.RowsAffected == 0 {
		return &user, false
	} else {
		return &user, true
	}
}

// 通过设备号，将游客更新为正式用户
func (dao *UserDAOImpl) VisitorToUser(device string, user *model.User) error {
	result := dao.db.Model(&user).Where("device_num = ?", device).Updates(&user)
	if result.Error != nil {

		return result.Error
	}
	return nil
}

// 创建用户或者游客
func (dao *UserDAOImpl) CreateUser(user *model.User) error {
	result := dao.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// todo 可以将这个函数改成更新用户的所有状态 姓名……等多种属性
// 通过设备号，更新游客的姓名
func (dao *UserDAOImpl) UpdateUserName(device string, username string) error {
	result := dao.db.Model(&model.User{}).Where("device_num = ?", device).Update("username", username)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (dao *UserDAOImpl) UpdateUser(user *model.User) error {
	result := dao.db.Table("users").Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通过邮箱，修改用户的密码
func (dao *UserDAOImpl) UpdatePassword(email string, password string) error {
	result := dao.db.Model(&model.User{}).Where("email = ?", email).Update("password", password)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通过邮箱，删除正式用户
func (dao *UserDAOImpl) DeleteUser(email string, user *model.User) error {
	result := dao.db.Model(&user).Where("email = ?", email).Unscoped().Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 更新用户的验证码
func (dao *UserDAOImpl) UpdateUserEmail(email *model.Email) error {
	result := dao.db.Save(&email)
	if result.Error != nil {
		return fmt.Errorf("更新验证码失败:%s", result.Error.Error())
	}
	return nil
}

// 创建用户邮箱验证码的内容
func (dao *UserDAOImpl) CreateUserEmail(email *model.Email) error {
	result := dao.db.Create(&email)
	if result.Error != nil {
		return fmt.Errorf("更新验证码失败%s", result.Error.Error())
	}
	return nil
}

// 通过用户id修改用户名
func (dao *UserDAOImpl) UpdateUserNameById(id uint, name string) error {
	result := dao.db.Model(&model.User{}).Where("id = ?", id).Update("username", name)
	if result.Error != nil {
		return fmt.Errorf("修改失败")
	}
	return nil
}

// 通过用户id修改邮箱
func (dao *UserDAOImpl) UpdateUserEmailById(id uint, email string) error {

	result := dao.db.Model(&model.User{}).Where("id = ?", id).Update("email", email)
	if result.Error != nil {
		return fmt.Errorf("修改邮箱失败")
	}
	return nil
}

// 随机一个用户
func (dao *UserDAOImpl) RandUser() (*model.User, error) {
	var user model.User
	result := dao.db.Model(&user).Order("rand()").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
