package middleware

import (
	"HarborArk/internal/i18n"
	"HarborArk/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": i18n.T(c, "token_missing"),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": i18n.T(c, "token_format_error"),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": i18n.T(c, "invalid_token"),
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("group_id", claims.GroupID)

		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(requiredLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID, exists := c.Get("group_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": i18n.T(c, "unauthorized"),
				"data":    nil,
			})
			c.Abort()
			return
		}

		userGroupID := groupID.(uint)
		if int(userGroupID) > requiredLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": i18n.T(c, "permission_denied"),
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}