package main

import (
	"database/sql"
	"fmt"
)

const (
	dbHost     = "localhost"
	dbPort     = 3306
	dbUser     = "root"
	dbPassword = "j72739906"
	dbName     = "conglong"
	secretKey  = "aComplexSecretKeyForSecurity"
)

var db *sql.DB

// 配置数据库连接
func InitDB() error {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		return err
	}
	return db.Ping()
}
