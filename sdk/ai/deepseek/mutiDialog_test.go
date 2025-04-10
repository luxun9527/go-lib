package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luxun9527/zlog"
	"io"
	"net/http"
	"strings"
	"testing"
)

// 定义结构体
type (
	DeepSeekRequest struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
		Stream   bool          `json:"stream"`
	}

	ChatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	DeepSeekStreamResponse struct {
		ID      string `json:"id"`
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}
)

func TestA(t *testing.T) {
	apiKey := "xxx" // 替换为你的 DeepSeek API Key
	apiURL := "https://api.deepseek.com/chat/completions"

	// 构造请求
	request := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "今天天气怎么样"},
			{Role: "system", Content: "如果你告诉我所在的城市或地区，我可以帮你查询最新的天气信息！"},
			{Role: "user", Content: "深圳"},
		},
		Stream: true,
	}

	// 发送请求并处理流式响应
	if err := streamDeepSeekResponse(apiKey, apiURL, request); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func streamDeepSeekResponse(apiKey, apiURL string, request DeepSeekRequest) error {
	// 序列化请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s\n%s", resp.Status, string(body))
	}
	var complateData string
	// 使用 bufio.Reader 逐行读取流式响应
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break // 流结束
			}
			return fmt.Errorf("read stream failed: %w", err)
		}

		// 清理数据并跳过空行
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		zlog.Infof("line %v", line)

		// 解析 JSON 数据
		var chunk DeepSeekStreamResponse
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			zlog.Infof("finish complate data %v", complateData)
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("unmarshal chunk failed: %w", err)
		}

		// 处理并打印内容
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			complateData += chunk.Choices[0].Delta.Content
			zlog.Infof("增量数据 %v", chunk.Choices[0].Delta.Content) // 实时输出
		}
	}

	return nil
}
