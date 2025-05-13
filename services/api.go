package services

import (
	"ImagetoEnglish/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ChatMessage struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type MessageContent struct {
	Type     string        `json:"type"`                // "text" 或 "image_url"
	Text     string        `json:"text,omitempty"`      // 当 type 为 "text" 时使用
	ImageURL *ImagePayload `json:"image_url,omitempty"` // 当 type 为 "image_url" 时使用
}

type ImagePayload struct {
	URL string `json:"url"`
}

// ChatRequest 表示请求体
type ChatRequest struct {
	Model           string        `json:"model"`
	Messages        []ChatMessage `json:"messages"`
	Response_format string        `json:"json_object"`
	Stream          bool          `json:"stream"`
	Temperature     float32       `json:"temperature"`
}

// ChatResponse 表示 Zhipu 的响应结构
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// CallAI 统一封装请求逻辑
func CallAI(modelKey string, messages []ChatMessage) (string, error) {
	aiConfig := config.Config.Sub("Ai")

	modelPath := "models." + modelKey
	modelConfig := aiConfig.Sub(modelPath)
	if modelConfig == nil {
		return "", fmt.Errorf("未找到模型配置: %s", modelKey)
	}

	apiKey := modelConfig.GetString("apiKey")
	if apiKey == "" {
		return "", fmt.Errorf("模型 %s 未配置 apiKey", modelKey)
	}

	baseURL := modelConfig.GetString("base_url")
	if baseURL == "" {
		return "", fmt.Errorf("模型 %s 未配置 base_url", modelKey)
	}

	modelName := modelConfig.GetString("model")
	if modelName == "" {
		return "", fmt.Errorf("模型 %s 未配置 model", modelKey)
	}

	temp := modelConfig.GetFloat64("temperature")
	if temp == 0 {
		temp = 0.95 // 给个默认值
	}

	// 构造请求
	reqBody := ChatRequest{
		Model:           modelName,
		Messages:        messages,
		Response_format: "json_object",
		Stream:          false,
		Temperature:     float32(temp),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		baseURL,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("AI 没有返回有效响应")
	}

	return chatResp.Choices[0].Message.Content, nil
}
