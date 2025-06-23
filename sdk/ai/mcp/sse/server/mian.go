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
		"demo-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add a test tool
	mcpServer.AddTool(mcp.NewTool(
		"getWeather",
		mcp.WithDescription("获取天气"),
		mcp.WithString("city", mcp.Description("城市名")),
		mcp.WithString("county", mcp.Description("区县名")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.GetArguments()
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
	log.Printf("init server at %v", "7878")
	if err := sseServer.Start(":7878"); err != nil {
		log.Printf("Error starting server: %v", err)
		return
	}

}
