package router

import (
	"HarborArk/cmd/docs"
	"HarborArk/config"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwagger 设置 Swagger 路由
func SetupSwagger(r *gin.Engine) {
	swaggerConfig := config.GetSwaggerConfig()

	if !swaggerConfig.Enabled {
		return
	}

	// 设置 Swagger 信息
	docs.SwaggerInfo.Title = swaggerConfig.Title
	docs.SwaggerInfo.Description = swaggerConfig.Description
	docs.SwaggerInfo.Version = swaggerConfig.Version
	docs.SwaggerInfo.Host = swaggerConfig.Host
	docs.SwaggerInfo.BasePath = swaggerConfig.BasePath
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger UI 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 文档重定向
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// API 信息端点
	r.GET("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title":       swaggerConfig.Title,
			"description": swaggerConfig.Description,
			"version":     swaggerConfig.Version,
			"docs_url":    "/swagger/index.html",
		})
	})
}
