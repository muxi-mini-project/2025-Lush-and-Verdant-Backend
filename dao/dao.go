package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewDB(addr string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值，超过此时间的查询将被记录
			LogLevel:      logger.Info, // 记录信息级别（Info、Warn、Error）
			Colorful:      true,        // 输出带颜色
		},
	)
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
