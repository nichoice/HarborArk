package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig 应用配置
type AppConfig struct {
	Server      ServerConfig      `mapstructure:"server"`
	Logger      LogConfig         `mapstructure:"logger"`
	Swagger     SwaggerConfig     `mapstructure:"swagger"`
	Metadata    MetadataConfig    `mapstructure:"metadata"`
	Audit       AuditConfig       `mapstructure:"audit"`
	FileManager FileManagerConfig `mapstructure:"fileManager"`
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

// BadgerConfig BadgerDB 配置
type BadgerConfig struct {
	Path string `mapstructure:"path"`
}

// MetadataConfig 元数据配置
type MetadataConfig struct {
	Badger BadgerConfig `mapstructure:"badger"`
}

// AuditConfig 审计配置
type AuditConfig struct {
	RetentionDays int    `mapstructure:"retentionDays"`
	ExportDir     string `mapstructure:"exportDir"`
}

// FileManagerConfig 文件管理器配置
type FileManagerConfig struct {
	RootDir               string   `mapstructure:"rootDir"`
	AllowedDirs           []string `mapstructure:"allowedDirs"`
	RestrictToAllowedDirs bool     `mapstructure:"restrictToAllowedDirs"`
	MaxDepth              int      `mapstructure:"maxDepth"`
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

// GetBadgerConfig 获取 Badger 配置
func GetBadgerConfig() BadgerConfig {
	if Config == nil {
		return BadgerConfig{
			Path: "data/badger",
		}
	}
	return Config.Metadata.Badger
}

// GetAuditConfig 获取审计配置
func GetAuditConfig() AuditConfig {
	if Config == nil {
		return AuditConfig{
			RetentionDays: 30,
			ExportDir:     "exports/audit",
		}
	}
	return Config.Audit
}

// GetFileManagerConfig 获取文件管理器配置
func GetFileManagerConfig() FileManagerConfig {
	if Config == nil {
		return FileManagerConfig{
			RootDir:               "/",
			AllowedDirs:           []string{"/"},
			RestrictToAllowedDirs: false,
			MaxDepth:              10,
		}
	}
	return Config.FileManager
}
