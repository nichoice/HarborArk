package config

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := "sqlite.db"

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zap.L().Error("连接数据库失败", zap.Error(err))
		panic(err)
	}

	zap.L().Info("连接数据库成功")

	DB.Logger = logger.Default.LogMode(logger.Info)

}
