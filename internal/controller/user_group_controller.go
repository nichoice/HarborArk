package controller

import (
	"HarborArk/internal/i18n"
	"HarborArk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserGroupController struct {
	userGroupService *service.UserGroupService
}

func NewUserGroupController() *UserGroupController {
	return &UserGroupController{
		userGroupService: &service.UserGroupService{},
	}
}

// CreateUserGroup 创建用户组
// @Summary 创建用户组
// @Description 创建新用户组
// @Tags 用户组管理
// @Accept json
// @Produce json
// @Param group body CreateUserGroupRequest true "用户组信息"
// @Success 200 {object} Response{data=model.UserGroup}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /user-groups [post]
func (c *UserGroupController) CreateUserGroup(ctx *gin.Context) {
	var req CreateUserGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_params"),
			Data:    nil,
		})
		return
	}

	group, err := c.userGroupService.CreateUserGroup(req.Name, req.Description, req.Level)
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
		Data:    group,
	})
}

// GetUserGroup 获取用户组详情
// @Summary 获取用户组详情
// @Description 根据ID获取用户组详情
// @Tags 用户组管理
// @Accept json
// @Produce json
// @Param id path int true "用户组ID"
// @Success 200 {object} Response{data=model.UserGroup}
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Security BearerAuth
// @Router /user-groups/{id} [get]
func (c *UserGroupController) GetUserGroup(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_group_id"),
			Data:    nil,
		})
		return
	}

	group, err := c.userGroupService.GetUserGroupByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: i18n.T(ctx, "group_not_found"),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "get_success"),
		Data:    group,
	})
}

// GetUserGroups 获取用户组列表
// @Summary 获取用户组列表
// @Description 分页获取用户组列表
// @Tags 用户组管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response{data=PageResponse}
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /user-groups [get]
func (c *UserGroupController) GetUserGroups(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	groups, total, err := c.userGroupService.GetUserGroups(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: i18n.T(ctx, "get_group_list_failed"),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: i18n.T(ctx, "get_success"),
		Data: PageResponse{
			List:     groups,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// UpdateUserGroup 更新用户组
// @Summary 更新用户组
// @Description 更新用户组信息
// @Tags 用户组管理
// @Accept json
// @Produce json
// @Param id path int true "用户组ID"
// @Param group body UpdateUserGroupRequest true "用户组信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /user-groups/{id} [put]
func (c *UserGroupController) UpdateUserGroup(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_group_id"),
			Data:    nil,
		})
		return
	}

	var req UpdateUserGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_params"),
			Data:    nil,
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Level != nil {
		updates["level"] = *req.Level
	}

	if err := c.userGroupService.UpdateUserGroup(uint(id), updates); err != nil {
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

// DeleteUserGroup 删除用户组
// @Summary 删除用户组
// @Description 删除用户组
// @Tags 用户组管理
// @Accept json
// @Produce json
// @Param id path int true "用户组ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Security BearerAuth
// @Router /user-groups/{id} [delete]
func (c *UserGroupController) DeleteUserGroup(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: i18n.T(ctx, "invalid_group_id"),
			Data:    nil,
		})
		return
	}

	if err := c.userGroupService.DeleteUserGroup(uint(id)); err != nil {
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
