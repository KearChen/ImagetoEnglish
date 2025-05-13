package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfig() {
	Config = viper.New()
	Config.SetConfigName("config")  // 配置文件名（不带扩展名）
	Config.SetConfigType("yaml")    // 配置文件类型
	Config.AddConfigPath("config/") // 配置文件所在目录

	if err := Config.ReadInConfig(); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 启用配置热更新
	Config.WatchConfig()
	Config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件被修改、已热更新\n")
	})
	fmt.Println("欢迎关注公众号：「陈同学的IPO之路」")
}
