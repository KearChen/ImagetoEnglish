package routers

import (
	"ImagetoEnglish/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 路由
	apiGroup := router.Group("/v1")
	{
		// 上传图片接口
		apiGroup.POST("/upload", controllers.UploadImage)
		// 获取分析结果接口
		apiGroup.POST("/analyze", controllers.AnalyzeImage)
	}

	// 前端资源静态路由
	router.Static("/static", "./static")
	// index.html
	router.LoadHTMLFiles("./templates/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	return router
}
