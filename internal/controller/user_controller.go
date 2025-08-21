package controller

import (
	"HarborArk/internal/i18n"
	"HarborArk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &service.UserService{},
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 200 {object} Response{data=model.User}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_params"),
			Data:    nil,
		})
		return
	}

	user, err := c.userService.CreateUser(req.Username, req.Password, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "create_success"),
		Data:    user,
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} Response{data=model.User}
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_user_id"),
			Data:    nil,
		})
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: i18n.T(ctx, "user_not_found"),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "get_success"),
		Data:    user,
	})
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response{data=PageResponse}
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := c.userService.GetUsers(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: i18n.T(ctx, "get_user_list_failed"),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "get_success"),
		Data: PageResponse{
			List:     users,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body UpdateUserRequest true "用户信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_user_id"),
			Data:    nil,
		})
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_params"),
			Data:    nil,
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Password != nil {
		updates["password"] = *req.Password
	}
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}

	if err := c.userService.UpdateUser(uint(id), updates); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "update_success"),
		Data:    nil,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_user_id"),
			Data:    nil,
		})
		return
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "delete_success"),
		Data:    nil,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} Response{data=LoginResponse}
// @Failure 400 {object} Response
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_params"),
			Data:    nil,
		})
		return
	}

	token, user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "login_success"),
		Data: LoginResponse{
			Token: token,
			User:  user,
		},
	})
}
