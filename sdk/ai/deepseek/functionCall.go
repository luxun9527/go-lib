package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiURL = "https://api.deepseek.com/chat/completions"
	apiKey = "sk-xxxxxx" // 替换为你的API密钥
	model  = "deepseek-chat"
)

// user assistant tool  发送
// 定义请求和响应结构体
type Message struct {
	Role       string  `json:"role"`
	Content    string  `json:"content"`
	ToolCalls  []*Tool `json:"tool_calls,omitempty"`
	ToolCallId string  `json:"tool_call_id,omitempty"`
}

type Tool struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}
type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Arguments   string     `json:"arguments,omitempty"`
	Parameters  *Parameter `json:"parameters,omitempty"`
}
type Parameter struct {
	Type       string               `json:"type"`
	Properties map[string]*Property `json:"properties,omitempty"`
	Required   []string             `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
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
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
}

// 获取当前时间的函数
func getCurrentTime() string {
	currentTime := time.Now()
	return fmt.Sprintf("现在是 %s", currentTime.Format("2006-01-02 15:04:05"))
}
func getLocation(location string) string {
	if strings.TrimSpace(location) != "" {
		return location
	}
	return "深圳"
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

var (
	tools = []Tool{{
		Type: "function",
		Function: Function{
			Name: "get_user_location",
			Description: `当服务需要用户位置信息才能继续时调用此函数。适用场景包括但不限于：
1. 用户请求天气查询但未提供具体位置
2. 用户查询"附近"的服务或场所(如美食、酒店等)
3. 需要基于位置提供个性化推荐时
4. 用户询问当前位置相关信息时
5. 如果用户提供了位置，则获取用户提供的位置信息
函数会返回用户当前位置或引导用户提供位置信息。`,
			Parameters: &Parameter{
				Type: "object",
				Properties: map[string]*Property{"location": {
					Type:        "string",
					Description: `如果用户提供了位置，则获取用户提供的位置信息`,
				}},
				Required: nil,
			},
		}},
	}
)

// user assistant tool
func main() {
	fmt.Println("欢迎使用智能助手！输入内容开始对话（输入 exit 退出）")

	// 定义工具
	c := Cli{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n用户输入: ")
		if !scanner.Scan() || scanner.Text() == "" {
			break
		}
		userInput := scanner.Text()
		// 创建消息数组
		c.doReq(userInput, "", false)
		// 发送请求

	}
}

type Cli struct {
	messages []Message
}

func (c *Cli) doReq(input, toolId string, isFuncInput bool) {
	if isFuncInput {
		c.messages = append(c.messages, Message{Role: "tool", ToolCallId: toolId, Content: input})
	} else {
		c.messages = append(c.messages, Message{Role: "user", Content: input})
	}

	request := ChatRequest{
		Messages:    c.messages,
		Model:       model,
		Tools:       tools,
		Temperature: 0.7,
	}
	log.Printf("%+v", request)
	response, err := sendChatRequest(request)
	if err != nil {
		fmt.Printf("请求出错: %v\n", err)
		return
	}

	if len(response.Choices) == 0 {
		fmt.Println("没有收到有效响应")
		return
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
			fmt.Printf("call func AI: %s\n", result)
		case "get_user_location":
			msg := Message{
				Role:      assistantMessage.Role,
				Content:   assistantMessage.Content,
				ToolCalls: assistantMessage.ToolCalls,
			}
			c.messages = append(c.messages, msg)
			result := getLocation(gjson.Get(toolCall.Function.Arguments, "location").String())
			fmt.Printf("call func AI: %s\n", result)
			c.doReq(result, toolCall.ID, true)

		default:
			fmt.Printf("AI: 未知函数调用: %s\n", functionName)
		}
	} else {
		// 直接输出模型响应
		fmt.Printf("AI: %s\n", assistantMessage.Content)
		c.messages = append(c.messages, Message{Role: "assistant", Content: assistantMessage.Content})
	}
}
