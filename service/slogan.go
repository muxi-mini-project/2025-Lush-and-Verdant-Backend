package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/dao"
	"fmt"
)

type SloganService interface {
	GetSlogan(string) (string, error)
	ChangeSlogan(uint, request.Slogan) error
}

type SloganServiceImpl struct {
	SloganDao dao.SloganDAO
	UserDao   dao.UserDAO
}

// 一个实例化对象
func NewSloganServiceImpl(sloganDao dao.SloganDAO, userDao dao.UserDAO) *SloganServiceImpl {
	return &SloganServiceImpl{
		SloganDao: sloganDao,
		UserDao:   userDao,
	}
}

func (ssr *SloganServiceImpl) GetSlogan(device string) (string, error) {
	slogans, err := ssr.SloganDao.GetAllSlogan()
	if err != nil {
		return "", err
	}

	// 如果没有可用的激励语，返回错误
	if len(slogans) == 0 {
		return "", fmt.Errorf("没有可用的激励语")
	}

	slogan, err := ssr.SloganDao.GetOneSlogan()
	if err != nil {
		return "", err
	}

	user, ok := ssr.UserDao.CheckUserByDevice(device)
	if !ok {
		return "", fmt.Errorf("找不到用户对应的设备号")
	}

	user.Slogan = slogan.Slogan
	err = ssr.UserDao.UpdateUser(user)
	if err != nil {
		return "", fmt.Errorf("更新激励语失败%s", err.Error())
	}

	return user.Slogan, nil
}

func (ssr *SloganServiceImpl) ChangeSlogan(id uint, newSlogan request.Slogan) error {
	user, err := ssr.UserDao.GetUserById(id)
	if err != nil {
		return fmt.Errorf("未找到相关用户%s", err.Error())
	}

	user.Slogan = newSlogan.Slogan
	err = ssr.UserDao.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
