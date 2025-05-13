package main

import (
	"ImagetoEnglish/config"
	"ImagetoEnglish/routers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// ------------------------- 配置初始化 -------------------------
	// 初始化全局配置文件
	config.InitConfig()

	// 获取并设置 Gin 运行模式（支持：debug、release）
	mode := config.Config.GetString("app.SetMode")
	gin.SetMode(mode)

	// ------------------------- 路由初始化 -------------------------
	// 初始化 Gin 路由
	app := routers.InitRouter()

	// ------------------------- 服务启动 -------------------------
	// 获取服务监听端口
	port := config.Config.GetInt("app.port")
	fmt.Printf("服务启动，端口：%d\n", port)

	// 启动 HTTP 服务
	err := app.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("服务启动失败: %v", err))
	}
}
