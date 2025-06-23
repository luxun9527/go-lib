package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type SseHandler struct {
}

func NewSseHandler() *SseHandler {
	return &SseHandler{}
}

// Serve 处理 SSE 连接
func (h *SseHandler) Serve(w http.ResponseWriter, r *http.Request) {
	// 设置 SSE 必需的 HTTP 头
	// for versions > v1.8.1, no need to add 3 lines below
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	// 为每个客户端创建一个 channel

	// 客户端断开时清理

	// 持续监听并推送事件

	for {
		select {
		case <-r.Context().Done():
			log.Println("finish")
			return
		default:
			time.Sleep(time.Second)
			// 发送事件数据
			message := fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC3339))

			fmt.Fprintf(w, "data: %s\n\n", message)
			f, ok := w.(http.Flusher)
			if ok {
				log.Println("ok", message)
				f.Flush()
			} else {
				log.Println("if not flusher")
			}
		}

	}
}
