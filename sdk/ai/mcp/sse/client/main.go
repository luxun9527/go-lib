package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"time"
)

func main() {
	client, err := client.NewSSEMCPClient("http://localhost:8888" + "/sse")
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the client
	if err := client.Start(ctx); err != nil {
		log.Println(err)
	}

	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	result, err := client.Initialize(ctx, initRequest)
	if err != nil {

		log.Printf("Initialize failed: %v", err)
	}

	if result.ServerInfo.Name != "test-server" {

	}

	// Test Ping
	if err := client.Ping(ctx); err != nil {
		log.Printf("Ping failed: %v", err)
	}

	// Test ListTools
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := client.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Printf("ListTools failed: %v", err)
	}
	for _, v := range tools.Tools {
		log.Printf("Tool: %v", v)
	}
}
