package main

import "context"

type Weather struct {
}

func (Weather) Name() string {
	return "get_weather"
}
func (Weather) Description() string {
	return "获取天气情况"
}
func (Weather) Call(ctx context.Context, input string) (string, error) {
	return "北京今天天气多云27C", nil
}
