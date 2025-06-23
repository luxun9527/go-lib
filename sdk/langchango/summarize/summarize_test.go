package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"testing"

	"github.com/ledongthuc/pdf"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/textsplitter"
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
func TestSummarizePDF(t *testing.T) {
	// 1. 加载 PDF 文件
	pdfPath := "test.pdf"
	text, err := extractTextFromPDF(pdfPath)
	if err != nil {
		log.Fatalf("PDF 解析失败: %v", err)
	}
	docs, err := documentloaders.NewText(strings.NewReader(text)).LoadAndSplit(context.Background(),
		textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(2000),
			textsplitter.WithChunkOverlap(200),
		),
	)

	if err != nil {
		log.Fatalf("文本分割失败: %v", err)
	}

	// 3. 初始化 OpenAI LLM
	llm, err := openai.New(
		openai.WithModel("qwen-turbo"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithToken("sk-xx"),
		openai.WithHTTPClient(loggingRoundTripper{}),
	)
	if err != nil {
		log.Println(err)
	}

	prompt := prompts.NewPromptTemplate(
		"总结以下内容，要求简洁准确：\n{{.context}}\n",
		[]string{"context"},
	)

	//	prompt1 := prompts.NewPromptTemplate(
	//		`你的任务是生成一个简洁的最终摘要
	//我们已提供截至某一点的现有摘要："{{.existing_answer}}"
	//我们有机会通过以下更多上下文来完善现有摘要（仅在需要时）。
	//"{{.context}}"
	//根据新的上下文，完善原始摘要
	//如果上下文无用，则返回原始摘要`,
	//		[]string{"existing_answer", "context"},
	//	)

	// 4. 使用总结链
	//	chain := chains.LoadRefineSummarization(chains.NewLLMChain(llm, prompt), chains.NewLLMChain(llm, prompt1))
	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(llm, prompt))

	chain := chains.NewMapReduceDocuments(chains.NewLLMChain(llm, prompt), combineChain)

	outputValues, err := chains.Call(context.Background(), chain, map[string]any{"input_documents": docs})
	if err != nil {
		log.Fatal(err)
	}
	out := outputValues["text"].(string)
	log.Println("out", out)
}

// 从 PDF 提取文本

// 核心函数：提取 PDF 所有页的文本
func extractTextFromPDF(filepath string) (string, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// 创建 PDF Reader
	reader, err := pdf.NewReader(file, fileInfo.Size())
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	// 逐页提取文本
	for pageNum := 1; pageNum <= reader.NumPage(); pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue // 跳过空白页
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("page %d error: %v", pageNum, err)
		}
		builder.WriteString(text + "\n") // 每页后加换行
	}

	return builder.String(), nil
}
