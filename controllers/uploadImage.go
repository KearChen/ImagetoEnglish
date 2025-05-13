package controllers

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed"})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	// 解码图片
	var img image.Image
	var decodeErr error
	if contentType == "image/png" {
		img, decodeErr = png.Decode(src)
	} else {
		img, _, decodeErr = image.Decode(src)
	}
	if decodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	// 创建静态目录
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		if err := os.Mkdir("./static", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	// 统一保存为 JPEG 格式
	randomName, err := generateRandomString(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate random filename"})
		return
	}
	newFilename := randomName + ".jpg"
	filePath := filepath.Join("./static", newFilename)

	// 压缩图片直到小于100KB
	var buf bytes.Buffer
	quality := 90
	for quality >= 10 {
		buf.Reset()
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compress image"})
			return
		}
		if buf.Len() < 100*1024 {
			break
		}
		quality -= 10
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save compressed image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("/static/%s", newFilename),
	})
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
