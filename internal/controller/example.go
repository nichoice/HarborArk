package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User 用户结构体
type User struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"张三"`
	Age  int    `json:"age" example:"25"`
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users := []User{
		{ID: 1, Name: "张三", Age: 25},
		{ID: 2, Name: "李四", Age: 30},
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    users,
		"message": "获取成功",
	})
}

// GetUser 根据ID获取用户
// @Summary 根据ID获取用户
// @Description 根据用户ID获取单个用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} User
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 模拟数据
	if id == 1 {
		user := User{ID: 1, Name: "张三", Age: 25}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    user,
			"message": "获取成功",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body User true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} map[string]interface{}
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 模拟创建用户
	user.ID = 3
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"data":    user,
		"message": "创建成功",
	})
}
