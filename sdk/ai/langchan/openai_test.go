package main

import (
	"context"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"log"
	"net/http"
	"net/http/httputil"
	"testing"
)

type loggerHttpCli struct {
}

func (loggerHttpCli) Do(req *http.Request) (*http.Response, error) {
	request, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	log.Printf("request:\n%s", string(request))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := httputil.DumpResponse(resp, true)
	log.Printf("response:\n%s", string(data))
	return resp, nil
}

func TestOpenai(t *testing.T) {

	llm, err := openai.New(
		openai.WithToken("sk-xxx"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithModel("qwen-turbo"),
		//openai.WithModel("qvq-plus"),
		openai.WithHTTPClient(loggerHttpCli{}),
	)
	if err != nil {
		log.Printf("Error creating LLM: %v", err)
		return
	}

	agentTools := []tools.Tool{Weather{}}

	agent := agents.NewOpenAIFunctionsAgent(llm,
		agentTools,
		agents.WithMaxIterations(3),
	)
	executor := agents.NewExecutor(agent)

	question := "北京今天天气怎么样？"
	answer, err := chains.Run(context.Background(), executor, question, chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		log.Printf("%s", chunk)
		return nil
	}))
	if err != nil {
		log.Printf("Error running agent: %v", err)
	}
	t.Logf("result: %s", answer)

}
