package dao

import (
	"2025-Lush-and-Verdant-Backend/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type MysqlDatabase struct {
	cfg *config.DatabaseConfig
}

func NewMysqlDatabase(cfg *config.DatabaseConfig) *MysqlDatabase {
	return &MysqlDatabase{
		cfg: cfg,
	}
}

func NewDB(md *MysqlDatabase) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值，超过此时间的查询将被记录
			LogLevel:      logger.Info, // 记录信息级别（Info、Warn、Error）
			Colorful:      true,        // 输出带颜色
		},
	)
	db, err := gorm.Open(mysql.Open(md.cfg.Addr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	// 获取数据库连接实例并设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(10)  // 设置最大空闲连接数

	return db, nil

}
