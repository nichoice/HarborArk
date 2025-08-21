package service

import (
	"HarborArk/config"
	"HarborArk/internal/model"
	"errors"

	"gorm.io/gorm"
)

type UserGroupService struct{}

// CreateUserGroup 创建用户组
func (s *UserGroupService) CreateUserGroup(name, description string, level int) (*model.UserGroup, error) {
	// 检查用户组名是否已存在
	var existingGroup model.UserGroup
	if err := config.DB.Where("name = ?", name).First(&existingGroup).Error; err == nil {
		return nil, errors.New("用户组名已存在")
	}

	// 创建用户组
	group := &model.UserGroup{
		Name:        name,
		Description: description,
		Level:       level,
	}

	if err := config.DB.Create(group).Error; err != nil {
		return nil, err
	}

	return group, nil
}

// GetUserGroupByID 根据ID获取用户组
func (s *UserGroupService) GetUserGroupByID(id uint) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := config.DB.Preload("Users").First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// GetUserGroups 获取用户组列表
func (s *UserGroupService) GetUserGroups(page, pageSize int) ([]model.UserGroup, int64, error) {
	var groups []model.UserGroup
	var total int64

	// 计算总数
	config.DB.Model(&model.UserGroup{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := config.DB.Preload("Users").Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// UpdateUserGroup 更新用户组
func (s *UserGroupService) UpdateUserGroup(id uint, updates map[string]interface{}) error {
	var group model.UserGroup
	if err := config.DB.First(&group, id).Error; err != nil {
		return err
	}

	return config.DB.Model(&group).Updates(updates).Error
}

// DeleteUserGroup 删除用户组
func (s *UserGroupService) DeleteUserGroup(id uint) error {
	var group model.UserGroup
	if err := config.DB.First(&group, id).Error; err != nil {
		return err
	}

	// 检查是否有用户使用此用户组
	var userCount int64
	config.DB.Model(&model.User{}).Where("group_id = ?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该用户组下还有用户，无法删除")
	}

	return config.DB.Delete(&group).Error
}

// InitDefaultUserGroups 初始化默认用户组
func (s *UserGroupService) InitDefaultUserGroups() error {
	groups := []model.UserGroup{
		{Name: "超级管理员", Description: "拥有系统所有权限", Level: model.SuperAdminGroup},
		{Name: "运维管理员", Description: "负责系统运维管理", Level: model.OpsAdminGroup},
		{Name: "审计管理员", Description: "负责系统审计工作", Level: model.AuditAdminGroup},
		{Name: "普通用户", Description: "普通系统用户", Level: model.NormalUserGroup},
	}

	for _, group := range groups {
		var existingGroup model.UserGroup
		if err := config.DB.Where("level = ?", group.Level).First(&existingGroup).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := config.DB.Create(&group).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}