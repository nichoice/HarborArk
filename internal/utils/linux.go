package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// CreateLinuxUser 创建Linux用户（nologin类型）
func CreateLinuxUser(username string) error {
	// 检查用户是否已存在
	if UserExists(username) {
		return fmt.Errorf("用户 %s 已存在", username)
	}

	// 创建用户，设置为nologin
	cmd := exec.Command("useradd", "-s", "/sbin/nologin", "-M", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("创建Linux用户失败: %v", err)
	}

	return nil
}

// DeleteLinuxUser 删除Linux用户
func DeleteLinuxUser(username string) error {
	if !UserExists(username) {
		return nil // 用户不存在，无需删除
	}

	cmd := exec.Command("userdel", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("删除Linux用户失败: %v", err)
	}

	return nil
}

// UserExists 检查Linux用户是否存在
func UserExists(username string) bool {
	cmd := exec.Command("id", username)
	err := cmd.Run()
	return err == nil
}

// EnableLinuxUser 启用Linux用户（修改shell）
func EnableLinuxUser(username string) error {
	cmd := exec.Command("usermod", "-s", "/bin/bash", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用Linux用户失败: %v", err)
	}
	return nil
}

// DisableLinuxUser 禁用Linux用户（设置为nologin）
func DisableLinuxUser(username string) error {
	cmd := exec.Command("usermod", "-s", "/sbin/nologin", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("禁用Linux用户失败: %v", err)
	}
	return nil
}

// GetLinuxUsers 获取系统用户列表
func GetLinuxUsers() ([]string, error) {
	cmd := exec.Command("cut", "-d:", "-f1", "/etc/passwd")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取系统用户列表失败: %v", err)
	}

	users := strings.Split(strings.TrimSpace(string(output)), "\n")
	return users, nil
}