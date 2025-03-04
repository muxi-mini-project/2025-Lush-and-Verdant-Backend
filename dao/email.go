package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type EmailDAO interface {
	GetEmail(addr string) (*model.Email, error)
	UpdateEmail(email *model.Email) error
}

// EmailCodeDAO 定义验证码存取接口
type EmailCodeDAO interface {
	SetEmailCode(email, code string, expiration time.Duration) error
	GetEmailCode(email string) (string, error)
	DeleteEmailCode(email string) error
}

type EmailDAOImpl struct {
	db *gorm.DB
}

// EmailCodeDAOImpl 实现 EmailCodeDAO
type EmailCodeDAOImpl struct {
	rdb *redis.Client
}

func NewEmailDAOImpl(db *gorm.DB) *EmailDAOImpl {
	return &EmailDAOImpl{
		db: db,
	}
}

// NewEmailCodeDAOImpl 创建新的 EmailCodeDAO
func NewEmailCodeDAOImpl(rdb *redis.Client) *EmailCodeDAOImpl {
	return &EmailCodeDAOImpl{rdb: rdb}
}

// 根据addr查询email
func (dao *EmailDAOImpl) GetEmail(addr string) (*model.Email, error) {
	var email model.Email
	//查询信息

	result := dao.db.Where("email=?", addr).First(&email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &email, nil
}

// 更新email
func (dao *EmailDAOImpl) UpdateEmail(email *model.Email) error {
	result := dao.db.Save(email)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// SetEmailCode 设置验证码，5分钟后过期
func (dao *EmailCodeDAOImpl) SetEmailCode(email, code string, expiration time.Duration) error {
	ctx := context.Background()
	err := dao.rdb.Set(ctx, email, code, expiration).Err()
	if err != nil {
		return fmt.Errorf("cannot set email code: %v", err)
	}
	return nil
}

// GetEmailCode 获取验证码
func (dao *EmailCodeDAOImpl) GetEmailCode(email string) (string, error) {
	ctx := context.Background()
	code, err := dao.rdb.Get(ctx, email).Result()
	if err == redis.Nil {
		return "", nil // 验证码不存在
	} else if err != nil {
		return "", fmt.Errorf("cannot get email code: %v", err)
	}
	return code, nil
}

// DeleteEmailCode 删除验证码
func (dao *EmailCodeDAOImpl) DeleteEmailCode(email string) error {
	ctx := context.Background()
	err := dao.rdb.Del(ctx, email).Err()
	if err != nil {
		return fmt.Errorf("cannot delete email code: %v", err)
	}
	return nil
}
