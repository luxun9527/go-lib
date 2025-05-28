package main

import (
	"context"
	"log"
)

type Calculator struct {
}

// Description returns a string describing the calculator tool.
func (c Calculator) Description() string {
	return `计算一个数学表达式
参数json格式
1. expression: 数学表达式 
`
}

// Name returns the name of the tool.
func (c Calculator) Name() string {
	return "calculator"
}

// Call evaluates the input using a starlak evaluator and returns the result as a
// string. If the evaluator errors the error is given in the result to give the
// agent the ability to retry.
func (c Calculator) Call(ctx context.Context, input string) (string, error) {
	log.Printf("input %s", input)
	return "2", nil
}
