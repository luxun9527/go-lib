package ali

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAli(t *testing.T) {
	var data = `{"messages":[{"content":"You are a helpful AI assistant.","role":"system"},{"content":"您是一个专业的天气助手，请用中文回答用户问题","role":"system"},{"content":"北京今天会不会下雨","role":"user"}],"model":"qwen-turbo","functions":[{"name":"get_current_weather","description":"查询指定城市的当前天气情况","parameters":{"properties":{"location":{"description":"城市名称，如'北京'","type":"string"}},"required":["location"],"type":"object"}}],"stream":false,"temperature":0.0}`
	req, _ := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", strings.NewReader(data))
	req.Header.Add("Authorization", "Bearer sk-xx")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(body))
}
