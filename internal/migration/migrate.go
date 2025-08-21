package migration

import (
	"HarborArk/config"
	"HarborArk/internal/model"
	"HarborArk/internal/service"
	"HarborArk/internal/utils"

	"go.uber.org/zap"
)

// AutoMigrate 自动迁移数据库
func AutoMigrate() error {
	// 自动迁移表结构
	if err := config.DB.AutoMigrate(&model.UserGroup{}, &model.User{}); err != nil {
		return err
	}

	zap.L().Info("数据库表结构迁移完成")

	// 初始化默认用户组
	userGroupService := &service.UserGroupService{}
	if err := userGroupService.InitDefaultUserGroups(); err != nil {
		return err
	}

	zap.L().Info("默认用户组初始化完成")

	// 创建默认超级管理员用户
	if err := createDefaultAdmin(); err != nil {
		zap.L().Warn("创建默认管理员失败", zap.Error(err))
	}

	return nil
}

// createDefaultAdmin 创建默认超级管理员
func createDefaultAdmin() error {
	userService := &service.UserService{}

	// 检查是否已存在admin用户
	_, err := userService.GetUserByUsername("admin")
	if err == nil {
		zap.L().Info("默认管理员用户已存在")
		return nil
	}

	// 创建默认管理员
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := &model.User{
		Username:  "admin",
		Password:  hashedPassword,
		GroupID:   model.SuperAdminGroup,
		IsEnabled: true,
	}

	if err := config.DB.Create(admin).Error; err != nil {
		return err
	}

	zap.L().Info("默认管理员用户创建成功", zap.String("username", "admin"), zap.String("password", "admin123"))
	return nil
}