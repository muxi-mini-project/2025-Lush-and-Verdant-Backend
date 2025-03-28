package service

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"2025-Lush-and-Verdant-Backend/model"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"time"
)

type ImageService interface {
	GetToken() (string, error)
	GetUserImage(user *model.User) (string, error)
	GetUserAllImage(user *model.User) ([]string, error)
	UpdateUserImage(image *model.UserImage) error
	GetGroupImage(group *model.Group) (string, error)
	GetGroupAllImage(group *model.Group) ([]string, error)
	UpdateGroupImage(image *model.GroupImage) error
}

type ImageServiceImpl struct {
	qny *config.QiNiuYunConfig
	Dao dao.ImageDAO
}

func NewImageServiceImpl(qny *config.QiNiuYunConfig, Dao dao.ImageDAO) *ImageServiceImpl {
	return &ImageServiceImpl{
		qny: qny,
		Dao: Dao,
	}
}

// 给前端传递token
func (isr *ImageServiceImpl) GetToken() (string, error) {
	accessKey := isr.qny.AccessKey
	secretKey := isr.qny.SecretKey
	bucket := isr.qny.BucketName
	//todo 这里不知道为什么莫名其妙丢了后面的两个
	fmt.Println(accessKey, " ", secretKey, " ", bucket, " ", isr.qny.BucketName, " ", isr.qny.DomainName)
	// mac是一个身份验证的工具
	mac := credentials.NewCredentials(accessKey, secretKey)

	// 生成一个创建策略，包含有效时间以及存储对象
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	// 获取上传凭证
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", err
	}
	return upToken, nil
}

// 通过UserId 获取user头像
func (isr *ImageServiceImpl) GetUserImage(user *model.User) (string, error) {
	url, err := isr.Dao.GetUserImage(user)
	if err != nil {
		return "", err
	}
	return url, nil
}

// 通过UserId 获取user历史头像
func (isr *ImageServiceImpl) GetUserAllImage(user *model.User) ([]string, error) {
	iamges, err := isr.Dao.GetUserAllImage(user)
	if err != nil {
		return nil, err
	}
	return iamges, nil
}

// 创建（用户更新图片）
func (isr *ImageServiceImpl) UpdateUserImage(image *model.UserImage) error {
	err := isr.Dao.CreateUserImage(image)
	if err != nil {
		return err
	}
	return nil
}

// group
// 通过groupId 获取group头像
func (isr *ImageServiceImpl) GetGroupImage(group *model.Group) (string, error) {
	url, err := isr.Dao.GetGroupImage(group)
	if err != nil {
		return "", err
	}
	return url, nil
}

// 通过groupId 获取group历史头像
func (isr *ImageServiceImpl) GetGroupAllImage(group *model.Group) ([]string, error) {
	images, err := isr.Dao.GetGroupAllImage(group)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// 创建（小组更新图片）
func (isr *ImageServiceImpl) UpdateGroupImage(image *model.GroupImage) error {
	//todo 检测权限
	err := isr.Dao.CreateGroupImage(image)
	if err != nil {
		return err
	}
	return nil
}
