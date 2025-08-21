package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig 应用配置
type AppConfig struct {
	Server  ServerConfig  `mapstructure:"server"`
	Logger  LogConfig     `mapstructure:"logger"`
	Swagger SwaggerConfig `mapstructure:"swagger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
}

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Title       string   `mapstructure:"title"`
	Description string   `mapstructure:"description"`
	Version     string   `mapstructure:"version"`
	Host        string   `mapstructure:"host"`
	BasePath    string   `mapstructure:"basePath"`
	Enabled     bool     `mapstructure:"enabled"`
	AutoUpdate  bool     `mapstructure:"autoUpdate"`
	OutputDir   string   `mapstructure:"outputDir"`
	MainApiFile string   `mapstructure:"mainApiFile"`
	Schemes     []string `mapstructure:"schemes"`
}

var Config *AppConfig

// Init 初始化配置
func Init() error {
	viper.SetConfigFile("config/settings-dev.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	Config = &AppConfig{}
	if err := viper.Unmarshal(Config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// GetLogConfig 获取日志配置
func GetLogConfig() LogConfig {
	if Config == nil {
		return LogConfig{
			Level:      "info",
			Encoding:   "json",
			Filename:   "logs/app.log",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 10,
		}
	}
	return Config.Logger
}

// GetServerConfig 获取服务器配置
func GetServerConfig() ServerConfig {
	if Config == nil {
		return ServerConfig{
			Port: "8080",
			Mode: "debug",
		}
	}
	return Config.Server
}

// GetSwaggerConfig 获取Swagger配置
func GetSwaggerConfig() SwaggerConfig {
	if Config == nil {
		return SwaggerConfig{
			Title:       "HarborArk API",
			Description: "HarborArk API Documentation",
			Version:     "1.0.0",
			Host:        "localhost:8080",
			BasePath:    "/api/v1",
			Enabled:     true,
			AutoUpdate:  false,
			OutputDir:   "cmd/docs",
			MainApiFile: "cmd/server.go",
			Schemes:     []string{"http", "https"},
		}
	}
	return Config.Swagger
}
