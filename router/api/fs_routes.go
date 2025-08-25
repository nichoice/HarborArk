package api

import (
	"HarborArk/internal/controller"
	"HarborArk/internal/middleware"
	"HarborArk/internal/model"

	"github.com/gin-gonic/gin"
)

// SetupFSRoutes 文件管理路由
func SetupFSRoutes(r *gin.Engine) {
	fs := controller.NewFSController()

	v1 := r.Group("/api/v1")
	group := v1.Group("/fs")
	group.Use(middleware.JWTAuth())
	{
		group.GET("/list", fs.List)
		group.GET("/metadata", fs.GetMetadata)
		group.PUT("/metadata", fs.UpdateMetadata)
		group.POST("/mkdir", middleware.RequireRole(model.SuperAdminGroup), fs.Mkdir)
		group.POST("/rename", middleware.RequireRole(model.SuperAdminGroup), fs.Rename)
		group.DELETE("/delete", middleware.RequireRole(model.SuperAdminGroup), fs.Delete)
		group.POST("/export-audit", middleware.RequireRole(model.SuperAdminGroup), fs.ExportAudit)
	}
}
