package controllers

import (
	"regexp"
	"strings"
)

func cleanAIResponse(raw string) string {
	// 1. 去除 markdown 代码块（```json ... ```）
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}
	// 修复 JSON 字符串中的非法换行符，将换行和前后的空白替换成空格
	reNewlineInString := regexp.MustCompile(`"([^"]*?)\s*\n\s*([^"]*?)"`)
	raw = reNewlineInString.ReplaceAllString(raw, `"$1 $2"`)
	// 2. 修复中文顿号：将 ["情感确认"、"安全感"] -> ["情感确认", "安全感"]
	// 正则匹配 ["xxx"、"yyy"、...]
	re := regexp.MustCompile(`\[\s*"[^"]*"(?:、\s*"[^"]*")+\s*\]`)
	raw = re.ReplaceAllStringFunc(raw, func(m string) string {
		// 去除首尾中括号后分割
		m = strings.TrimPrefix(m, "[")
		m = strings.TrimSuffix(m, "]")
		parts := strings.Split(m, "、")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return "[ " + strings.Join(parts, ", ") + " ]"
	})

	return raw
}
