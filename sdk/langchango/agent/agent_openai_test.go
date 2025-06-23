package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"log"
	"net/http"
	"net/http/httputil"
	"testing"
)

// 自定义 HTTP 客户端（启用请求日志）
type loggingRoundTripper struct {
}

func (lrt loggingRoundTripper) Do(req *http.Request) (*http.Response, error) {
	// 打印请求详情
	requestDump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("HTTP Request:\n%s\n", string(requestDump))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// 打印响应详情
	responseDump, _ := httputil.DumpResponse(resp, true)
	log.Printf("HTTP Response:\n%s\n", string(responseDump))
	return resp, nil
}
func TestOpenaiFunctionAgent(t *testing.T) {
	llm, err := openai.New(
		openai.WithModel("qwen-turbo"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithToken("sk-xx"),
		openai.WithHTTPClient(loggingRoundTripper{}),
	)
	if err != nil {
		log.Println(err)
	}
	//search, err := serpapi.New()
	//if err != nil {
	//	return err
	//}
	agentTools := []tools.Tool{
		Weather{},
	}
	agent := agents.NewOneShotAgent(llm,
		agentTools,
		agents.WithMaxIterations(1),
	)

	//agents.NewConversationalAgent()
	executor := agents.NewExecutor(agent)
	question := "北京？"
	answer, err := chains.Run(context.Background(), executor, question)
	fmt.Println(answer)

}
func TestOpenaiFunctionAgent1(t *testing.T) {
	llm, err := openai.New(
		openai.WithModel("qwen-turbo"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithToken("sk-xxxxxx"),
		openai.WithHTTPClient(loggingRoundTripper{}),
	)
	if err != nil {
		log.Println(err)
	}
	//search, err := serpapi.New()
	//if err != nil {
	//	return err
	//}
	agentTools := []tools.Tool{
		Weather{},
	}

	agent := agents.NewOpenAIFunctionsAgent(llm,
		agentTools,
		agents.WithMaxIterations(2),
	)

	agentAli := &OpenAIFunctionsAgentAli{agent}

	//agents.NewConversationalAgent()
	executor := agents.NewExecutor(agentAli)

	question := "今天北京天气如何？"
	answer, err := chains.Run(context.Background(), executor, question)
	log.Printf("answer: %s", answer)

	question = "上海？"
	answer, err = chains.Run(context.Background(), executor, question)
	log.Printf("answer: %s", answer)

}
func TestOpenaiFunctionConversational(t *testing.T) {
	llm, err := openai.New(
		openai.WithModel("qwen-turbo"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithToken("sk-xx"),
		openai.WithHTTPClient(loggingRoundTripper{}),
	)
	if err != nil {
		log.Println(err)
	}
	//search, err := serpapi.New()
	//if err != nil {
	//	return err
	//}
	agentTools := []tools.Tool{
		Weather{},
	}

	agent := agents.NewConversationalAgent(llm,
		agentTools,
		agents.WithMaxIterations(2),
	)

	//agents.NewConversationalAgent()
	executor := agents.NewExecutor(agent)

	question := "今天北京天气如何？"
	answer, err := chains.Run(context.Background(), executor, question)
	log.Printf("answer: %s", answer)
	answer, err = chains.Run(context.Background(), executor, question)
	log.Printf("answer: %s", answer)
}
