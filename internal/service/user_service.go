package service

import (
	"HarborArk/config"
	"HarborArk/internal/model"
	"HarborArk/internal/utils"
	"errors"
	"time"
)

type UserService struct{}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, password string, groupID uint) (*model.User, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username:  username,
		Password:  hashedPassword,
		GroupID:   groupID,
		IsEnabled: true,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// 如果是普通用户，同步创建Linux用户
	if groupID == model.NormalUserGroup {
		if err := utils.CreateLinuxUser(username); err != nil {
			// 如果Linux用户创建失败，记录日志但不回滚数据库操作
			// 可以考虑添加日志记录
		}
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := config.DB.Preload("Group").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := config.DB.Preload("Group").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 计算总数
	config.DB.Model(&model.User{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := config.DB.Preload("Group").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return err
	}

	// 如果更新密码，需要加密
	if password, ok := updates["password"]; ok {
		hashedPassword, err := utils.HashPassword(password.(string))
		if err != nil {
			return err
		}
		updates["password"] = hashedPassword
	}

	// 如果更新启用状态且是普通用户，同步更新Linux用户状态
	if isEnabled, ok := updates["is_enabled"]; ok && user.GroupID == model.NormalUserGroup {
		if isEnabled.(bool) {
			utils.EnableLinuxUser(user.Username)
		} else {
			utils.DisableLinuxUser(user.Username)
		}
	}

	return config.DB.Model(&user).Updates(updates).Error
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return err
	}

	// 如果是普通用户，同步删除Linux用户
	if user.GroupID == model.NormalUserGroup {
		utils.DeleteLinuxUser(user.Username)
	}

	return config.DB.Delete(&user).Error
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, *model.User, error) {
	// 获取用户
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户是否启用
	if !user.IsEnabled {
		return "", nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	config.DB.Model(user).Update("last_login_at", &now)

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.GroupID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}