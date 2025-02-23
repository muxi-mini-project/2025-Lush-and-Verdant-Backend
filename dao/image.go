package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"fmt"
	"gorm.io/gorm"
)

type ImageDAO interface {
	CreateUserImage(image *model.UserImage) error
	GetUserImage(user *model.User) (string, error)
	GetUserAllImage(user *model.User) ([]string, error)
	CreateGroupImage(image *model.GroupImage) error
	GetGroupImage(group *model.Group) (string, error)
	GetGroupAllImage(group *model.Group) ([]string, error)
}
type ImageDAOImpl struct {
	db *gorm.DB
}

func NewImageDAO(db *gorm.DB) *ImageDAOImpl {
	return &ImageDAOImpl{
		db: db,
	}
}

//通过userId获取image的最新url

func (dao *ImageDAOImpl) GetUserImage(user *model.User) (string, error) {
	dao.db.Preload("Images").Find(user)
	if len(user.Images) == 0 {
		return "", fmt.Errorf("没有头像")
	}
	return user.Images[len(user.Images)-1].Url, nil
}

// 通过userId获取image的所有url
func (dao *ImageDAOImpl) GetUserAllImage(user *model.User) ([]string, error) {
	dao.db.Preload("Images").Find(user)
	if len(user.Images) == 0 {
		return nil, fmt.Errorf("没有头像")
	}
	var images []string
	for _, image := range user.Images {
		images = append(images, image.Url)
	}
	return images, nil
}

// 创建
func (dao *ImageDAOImpl) CreateUserImage(image *model.UserImage) error {
	result := dao.db.Create(image)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通过groupId来获取头像
func (dao *ImageDAOImpl) GetGroupImage(group *model.Group) (string, error) {
	dao.db.Preload("Images").Find(group)
	if len(group.Images) == 0 {
		return "", fmt.Errorf("没有头像")
	}
	return group.Images[len(group.Images)-1].Url, nil
}

// 通过groupId来获取历史头像
func (dao *ImageDAOImpl) GetGroupAllImage(group *model.Group) ([]string, error) {
	dao.db.Preload("Images").Find(group)
	if len(group.Images) == 0 {
		return nil, fmt.Errorf("没有头像")
	}
	var images []string
	for _, image := range group.Images {
		images = append(images, image.Url)
	}
	return images, nil
}

// 创建
func (dao *ImageDAOImpl) CreateGroupImage(image *model.GroupImage) error {
	result := dao.db.Create(image)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
