package main

import (
	"database/sql"
	"log"
)

func init() {
	dsn := "root:Lu03150079@tcp(116.62.179.155:3306)/users"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//检查数据库连接
	if err := db.Ping(); err != nil {
		log.Println(err.Error())
		return
	}
}
