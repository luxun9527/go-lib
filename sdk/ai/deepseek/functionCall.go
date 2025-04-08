package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiURL = "https://api.deepseek.com/chat/completions"
	apiKey = "sk-xxxx" // 替换为你的API密钥
	model  = "deepseek-chat"
)

// 定义请求和响应结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"function"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type AssistantMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      AssistantMessage `json:"message"`
		FinishReason string           `json:"finish_reason"`
	} `json:"choices"`
}

// 获取当前时间的函数
func getCurrentTime() string {
	currentTime := time.Now()
	return fmt.Sprintf("现在是 %s", currentTime.Format("2006-01-02 15:04:05"))
}

func sendChatRequest(request ChatRequest) (*ChatResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &response, nil
}

func main() {
	fmt.Println("欢迎使用智能助手！输入内容开始对话（输入 exit 退出）")

	// 定义工具
	tools := []Tool{
		{
			Type: "function",
			Function: struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}{
				Name:        "get_current_time",
				Description: "当你想知道现在的时间时非常有用。",
			},
		},
	}

	for {
		fmt.Print("\n用户输入: ")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("再见！")
			os.Exit(0)
		}

		// 创建消息数组
		messages := []Message{
			{Role: "system", Content: ""},
			{Role: "user", Content: userInput},
		}

		// 发送请求
		request := ChatRequest{
			Messages:    messages,
			Model:       model,
			Tools:       tools,
			Temperature: 0.7,
		}

		response, err := sendChatRequest(request)
		if err != nil {
			fmt.Printf("请求出错: %v\n", err)
			continue
		}

		if len(response.Choices) == 0 {
			fmt.Println("没有收到有效响应")
			continue
		}

		assistantMessage := response.Choices[0].Message

		// 检查是否有工具调用
		if len(assistantMessage.ToolCalls) > 0 {
			toolCall := assistantMessage.ToolCalls[0]
			functionName := toolCall.Function.Name

			// 处理函数调用
			switch functionName {
			case "get_current_time":
				result := getCurrentTime()
				fmt.Printf("AI: %s\n", result)
			default:
				fmt.Printf("AI: 未知函数调用: %s\n", functionName)
			}
		} else {
			// 直接输出模型响应
			fmt.Printf("AI: %s\n", assistantMessage.Content)
		}
	}
}
