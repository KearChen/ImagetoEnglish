package controllers

import (
	"ImagetoEnglish/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ChineseName string `json:"中文名称"`
	EnglishName string `json:"英文名称"`
	EnglishDesc string `json:"英文描述"`
}

func AnalyzeImage(c *gin.Context) {
	// 获取图片 URL
	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 构造提示词
	systemPrompt := `你是一个图像识别和语言学习助手。用户上传一张图片，你的任务是：
    判断图中是什么物体（例如：苹果、书包、椅子、眼镜、狗、猫等）；
    给出该物体的标准中文名称和对应英文名称；
    提供一句最简单的英文描述。`
	exampleFormat := `请严格按照如下格式仅输出JSON，不要输出代码或其他信息，JSON字段使用顿号【、】区隔：{
"中文名称": "",
"英文名称": "",
"英文描述": ""
}`

	messages := []services.ChatMessage{
		{
			Role: "system",
			Content: []services.MessageContent{
				{Type: "text", Text: systemPrompt},
			},
		},
		{
			Role: "user",
			Content: []services.MessageContent{
				{Type: "text", Text: exampleFormat},
			},
		},
		{
			Role: "user",
			Content: []services.MessageContent{
				{Type: "text", Text: "请识别图片并返回中英文名称与描述"},
				{Type: "image_url", ImageURL: &services.ImagePayload{URL: req.ImageURL}},
			},
		},
	}

	respText, err := services.CallAI("glm-4v-flash", messages)
	if err != nil {
		log.Printf("请求失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 响应失败"})
		return
	}

	sanitized := cleanAIResponse(respText)

	var parsed Response
	if err := json.Unmarshal([]byte(sanitized), &parsed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 返回格式错误"})
		return
	}

	analysisResult := map[string]interface{}{
		"中文名称": parsed.ChineseName,
		"英文名称": parsed.EnglishName,
		"英文描述": parsed.EnglishDesc,
	}

	c.JSON(http.StatusOK, analysisResult)
}
