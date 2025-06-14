package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"log"
)

type Weather struct {
	CallbacksHandler callbacks.Handler
}

// Description returns a string describing the calculator tool.
func (c Weather) Description() string {
	m := map[string]interface{}{
		"location": map[string]string{
			"type":        "string",
			"description": "用户的位置，可以是城市名，也可以为空",
		},
		"datetimeType": map[string]string{
			"type": "string",
			"description": `用户要查询天气的的时间类型 
now:实时天气当用户查询今天或现在天气
day:用户查询的每日天气预报，如明天天气,后天天气
hour:用户查询24小时，逐小时天气预报，如一小时会不会下雨`,
		},
	}
	d, _ := json.Marshal(m)
	return fmt.Sprintf(`获取实时天气数据,24小时内每小时的天气预报信息，以及7天内每天的天气预报，不支持查询历史天气。支持全球城市精确查询，支持不指定位置。
参数: %s
`, string(d))
}

// Name returns the name of the tool.
func (c Weather) Name() string {
	return "getWeather"
}

// Call evaluates the input using a starlak evaluator and returns the result as a
// string. If the evaluator errors the error is given in the result to give the
// agent the ability to retry.
func (c Weather) Call(ctx context.Context, input string) (string, error) {

	log.Printf("input %s", input)
	return "北京天气是23摄氏度，多云", nil
}

func (c Weather) Parameters() map[string]any {
	m := map[string]interface{}{
		"location": map[string]string{
			"type":        "string",
			"description": "用户的位置，可以是城市名，也可以为空",
		},
		"datetimeType": map[string]string{
			"type": "string",
			"description": `用户要查询天气的的时间类型 
now:实时天气当用户查询今天或现在天气
day:用户查询的每日天气预报，如明天天气,后天天气
hour:用户查询24小时，逐小时天气预报，如一小时会不会下雨`,
		},
	}
	return m
}
