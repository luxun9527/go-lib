package main

import (
	"context"
	"github.com/luxun9527/zlog"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"testing"
	"time"
)

func TestClient1(t *testing.T) {
	client, err := client.NewSSEMCPClient("http://localhost:7878" + "/sse")
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

	_, err = client.Initialize(ctx, initRequest)
	if err != nil {

		log.Printf("Initialize failed: %v", err)
	}

	// Test Ping
	if err := client.Ping(ctx); err != nil {
		log.Printf("Ping failed: %v", err)
	}

	// Test ListTools
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := client.ListTools(ctx, toolsRequest)
	if err != nil {
		zlog.Infof("listTools failed: %+v", err)
		return
	}
	for _, v := range tools.Tools {
		zlog.Infof("tool: %+v", v)
		tool, err := client.CallTool(ctx, mcp.CallToolRequest{
			Request: mcp.Request{
				Method: v.Name,
				Params: mcp.RequestParams{Meta: &mcp.Meta{
					ProgressToken:    nil,
					AdditionalFields: nil,
				}},
			},
			Params: mcp.CallToolParams{
				Name: v.Name,
				Arguments: map[string]any{
					"city": "北京",
				},
				Meta: nil,
			},
		})
		if err != nil {
			zlog.Errorf("callTool failed: %+v", err)
		}
		zlog.Infof("tool: %+v", tool)
	}

}
func TestClientMap(t *testing.T) {
	client, err := client.NewSSEMCPClient("https://mcp.amap.com/sse?key=xxxx")
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
	//initRequest.Params.ClientInfo = mcp.Implementation{
	//	Name:    "test-client",
	//	Version: "1.0.0",
	//}

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
		zlog.Infof("listTools failed: %+v", err)
		return
	}
	zlog.Infof("tools: %+v", tools)
	for _, v := range tools.Tools {
		zlog.Infof("tool: %+v", v)
	}
	for _, v := range tools.Tools {
		zlog.Infof("tool: %+v", v)
		tool, err := client.CallTool(ctx, mcp.CallToolRequest{
			Request: mcp.Request{
				Method: v.Name,
				Params: mcp.RequestParams{Meta: &mcp.Meta{
					ProgressToken:    nil,
					AdditionalFields: nil,
				}},
			},
			Params: mcp.CallToolParams{
				Name:      v.Name,
				Arguments: nil,
				Meta:      nil,
			},
		})
		if err != nil {
			zlog.Errorf("callTool failed: %+v", err)
		}
		zlog.Infof("tool: %+v", tool)
	}

}
