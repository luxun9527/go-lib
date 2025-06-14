package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

func main() {
	// Create MCP server with capabilities
	mcpServer := server.NewMCPServer(
		"test-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	// Add a test tool
	mcpServer.AddTool(mcp.NewTool(
		"test-tool",
		mcp.WithDescription("测试工具"),
		mcp.WithObject("city", mcp.Properties(map[string]any{
			"city": map[string]any{"type": "string", "description": "城市名称"},
		})),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.GetArguments()
		a := request.GetString("city", "f")
		log.Printf("city: %v", a)
		log.Printf("Arguments: %v", arguments)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "北京今天多云。",
				},
			},
		}, nil
	})
	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://%s", "localhost:7878")))
	err := sseServer.Start(":7878")
	if err != nil {
		log.Printf("Error starting server: %v", err)
		return
	}
}
