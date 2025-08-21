package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type (
	// Runmode 开发模式
	Runmode string
)

const (
	RunmodeDev  Runmode = "dev"
	RunmodeProd Runmode = "prod"
	RunmodeTest Runmode = "test"
)

func Setup(configFilePath, configFileType string) {
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType(configFileType)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败，请检查 %s 配置文件是否存在: %v", configFilePath, err))
	}
	viper.AutomaticEnv()                                   // 自动读取环境变量
	viper.SetEnvPrefix("harborark")                        // 设置环境变量前缀
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 替换点

	watchConfigChanges()

	fmt.Printf("\nConfiguration file loaded successfully: %s\n", configFilePath)
}

// WriteConfig 写入配置文件
func WriteConfig(filename string) {
	viper.WriteConfigAs(filename)
}

func watchConfigChanges() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
}
