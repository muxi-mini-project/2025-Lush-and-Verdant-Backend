package dao

import (
	"2025-Lush-and-Verdant-Backend/model"
	"gorm.io/gorm"
)

type EmailDAO interface {
	GetEmail(addr string) (*model.Email, error)
	UpdateEmail(email *model.Email) error
}

type EmailDAOImpl struct {
	db *gorm.DB
}

func NewEmailDAOImpl(db *gorm.DB) *EmailDAOImpl {
	return &EmailDAOImpl{
		db: db,
	}

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
