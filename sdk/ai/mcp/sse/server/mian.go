package main

import (
	"context"
	"errors"
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
		mcp.WithDescription("Test tool"),
		mcp.WithString("parameter-1", mcp.Description("A string tool parameter")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Input parameter: " + request.Params.Arguments["parameter-1"].(string),
				},
			},
		}, nil
	})
	sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://%s", "localhost:8888")))
	err := sseServer.Start(":8888")
	if err != nil {
		log.Printf("Error starting server: %v", err)
		return
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
