package dao

import (
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/model"
	"log"
)

var dsn = config.Dsn

func Creat() {
	db := NewDB(dsn)
	err := db.AutoMigrate(&model.Email{})
	if err != nil {
		log.Println(err)
		return
	}
}
