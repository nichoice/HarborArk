package controller

import "HarborArk/internal/model"

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 分页响应结构
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 用户相关请求结构
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
	GroupID  uint   `json:"group_id" binding:"required" example:"1"`
}

type UpdateUserRequest struct {
	Password  *string `json:"password,omitempty" example:"newpassword"`
	GroupID   *uint   `json:"group_id,omitempty" example:"2"`
	IsEnabled *bool   `json:"is_enabled,omitempty" example:"true"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *model.User  `json:"user"`
}

// 用户组相关请求结构
type CreateUserGroupRequest struct {
	Name        string `json:"name" binding:"required" example:"管理员"`
	Description string `json:"description" example:"系统管理员组"`
	Level       int    `json:"level" binding:"required" example:"1"`
}

type UpdateUserGroupRequest struct {
	Name        *string `json:"name,omitempty" example:"新管理员"`
	Description *string `json:"description,omitempty" example:"新的描述"`
	Level       *int    `json:"level,omitempty" example:"2"`
}