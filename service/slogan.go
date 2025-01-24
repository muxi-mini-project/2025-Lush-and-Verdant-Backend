package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type SloganService struct {
	db *gorm.DB
}

func NewSloganService(db *gorm.DB) *SloganService {
	return &SloganService{
		db: db,
	}
}

func (ssr *SloganService) GetSlogan(device string) error {

	// 查找所有的激励语
	var slogans []model.Slogan
	if err := ssr.db.Find(&slogans).Error; err != nil {
		return fmt.Errorf("无法获取激励语")
	}
	// 如果没有可用的激励语，返回错误
	if len(slogans) == 0 {
		return fmt.Errorf("没有可用的激励语")
	}

	var slogan model.Slogan
	ssr.db.Table("slogans").Order("RAND()").First(&slogan) // 使用Order随机排序，选第一条slogan

	var user model.User
	err := ssr.db.Table("users").Where("device_num = ?", device).First(&user)
	if err.Error != nil {
		return fmt.Errorf("找不到用户对应的设备号")
	}

	user.Slogan = slogan.Slogan
	err = ssr.db.Table("users").Save(&user)
	if err.Error != nil {
		return fmt.Errorf("更新激励语失败")
	}
	return nil
}

func (ssr *SloganService) ChangeSlogan(id uint, newSlogan request.Slogan) error {

	var user model.User
	err := ssr.db.Table("users").Where("id = ?", id).First(&user).Error
	if err.Error != nil {
		return fmt.Errorf("未找到相关用户")
	}
	user.Slogan = newSlogan.Slogan
	err = ssr.db.Table("users").Save(&user).Error
	if err.Error != nil {
		return err
	}
	return nil
}
