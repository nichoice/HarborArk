package model

import (
	"time"

	"gorm.io/gorm"
)

// UserGroup 用户组模型
type UserGroup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;comment:用户组名称"`
	Description string         `json:"description" gorm:"comment:用户组描述"`
	Level       int            `json:"level" gorm:"not null;comment:权限级别"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Users       []User         `json:"users,omitempty" gorm:"foreignKey:GroupID"`
}

// User 用户模型
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;not null;comment:用户名"`
	Password    string         `json:"-" gorm:"not null;comment:密码"`
	GroupID     uint           `json:"group_id" gorm:"not null;comment:用户组ID"`
	Group       UserGroup      `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	LastLoginAt *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserGroupType 用户组类型常量
const (
	SuperAdminGroup = iota + 1 // 超级管理员
	OpsAdminGroup              // 运维管理员
	AuditAdminGroup            // 审计管理员
	NormalUserGroup            // 普通用户
)

// GetUserGroupName 获取用户组名称
func GetUserGroupName(level int) string {
	switch level {
	case SuperAdminGroup:
		return "超级管理员"
	case OpsAdminGroup:
		return "运维管理员"
	case AuditAdminGroup:
		return "审计管理员"
	case NormalUserGroup:
		return "普通用户"
	default:
		return "未知"
	}
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (UserGroup) TableName() string {
	return "user_groups"
}
