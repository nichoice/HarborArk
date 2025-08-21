package api

import (
	"HarborArk/internal/controller"
	"HarborArk/internal/middleware"
	"HarborArk/internal/model"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(r *gin.Engine) {
	userController := controller.NewUserController()
	userGroupController := controller.NewUserGroupController()

	// API版本组
	v1 := r.Group("/api/v1")

	// 认证路由（无需JWT验证）
	auth := v1.Group("/auth")
	{
		auth.POST("/login", userController.Login)
	}

	// 用户路由（需要JWT验证）
	users := v1.Group("/users")
	users.Use(middleware.JWTAuth())
	{
		users.GET("", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.POST("", middleware.RequireRole(model.SuperAdminGroup), userController.CreateUser)
		users.PUT("/:id", middleware.RequireRole(model.SuperAdminGroup), userController.UpdateUser)
		users.DELETE("/:id", middleware.RequireRole(model.SuperAdminGroup), userController.DeleteUser)
	}

	// 用户组路由（需要JWT验证）
	userGroups := v1.Group("/user-groups")
	userGroups.Use(middleware.JWTAuth())
	{
		userGroups.GET("", userGroupController.GetUserGroups)
		userGroups.GET("/:id", userGroupController.GetUserGroup)
		userGroups.POST("", middleware.RequireRole(model.SuperAdminGroup), userGroupController.CreateUserGroup)
		userGroups.PUT("/:id", middleware.RequireRole(model.SuperAdminGroup), userGroupController.UpdateUserGroup)
		userGroups.DELETE("/:id", middleware.RequireRole(model.SuperAdminGroup), userGroupController.DeleteUserGroup)
	}
}
