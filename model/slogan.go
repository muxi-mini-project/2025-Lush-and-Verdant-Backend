package model

type Slogan struct {
	ID     int    `gorm:"primary_key" `
	Slogan string `gorm:"varchar(255)" `
}
