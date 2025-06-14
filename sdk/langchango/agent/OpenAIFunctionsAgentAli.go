package main

import (
	"context"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

const agentScratchpad = "agent_scratchpad"

type OpenAIFunctionsAgentAli struct {
	*agents.OpenAIFunctionsAgent
}

func (agentAli *OpenAIFunctionsAgentAli) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for key, value := range inputs {
		fullInputs[key] = value
	}
	fullInputs[agentScratchpad] = agentAli.constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if agentAli.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			agentAli.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}

	prompt, err := agentAli.Prompt.FormatPrompt(fullInputs)
	if err != nil {
		return nil, nil, err
	}

	mcList := make([]llms.MessageContent, len(prompt.Messages()))
	for i, msg := range prompt.Messages() {
		role := msg.GetType()
		text := msg.GetContent()

		var mc llms.MessageContent

		switch p := msg.(type) {
		case llms.ToolChatMessage:
			mc = llms.MessageContent{
				Role: role,
				Parts: []llms.ContentPart{llms.ToolCallResponse{
					ToolCallID: p.ID,
					Content:    p.Content,
				}},
			}

		case llms.AIChatMessage:
			mc = llms.MessageContent{
				Role: role,
				Parts: []llms.ContentPart{
					llms.ToolCall{
						ID:           p.ToolCalls[0].ID,
						Type:         p.ToolCalls[0].Type,
						FunctionCall: p.ToolCalls[0].FunctionCall,
					},
				},
			}
		default:
			mc = llms.MessageContent{
				Role:  role,
				Parts: []llms.ContentPart{llms.TextContent{Text: text}},
			}
		}
		mcList[i] = mc
	}

	result, err := agentAli.LLM.GenerateContent(ctx, mcList,
		llms.WithTools(agentAli.functions()), llms.WithStreamingFunc(stream))
	if err != nil {
		return nil, nil, err
	}

	return agentAli.ParseOutput(result)
}

func (agentAli *OpenAIFunctionsAgentAli) constructScratchPad(steps []schema.AgentStep) []llms.ChatMessage {
	if len(steps) == 0 {
		return nil
	}
	messages := make([]llms.ChatMessage, 0)
	for _, step := range steps {
		messages = append(messages, llms.AIChatMessage{
			ToolCalls: []llms.ToolCall{
				{
					ID:           step.Action.ToolID,
					Type:         "function",
					FunctionCall: &llms.FunctionCall{Name: step.Action.Tool, Arguments: step.Action.ToolInput},
				},
			},
		})
		messages = append(messages, llms.ToolChatMessage{
			ID:      step.Action.ToolID,
			Content: step.Observation,
		})
	}

	return messages
}

func (agentAli *OpenAIFunctionsAgentAli) functions() []llms.Tool {
	res := make([]llms.Tool, 0)
	for _, tool := range agentAli.Tools {
		t, ok := tool.(OpenaiTool)
		if ok {
			res = append(res, llms.Tool{
				Type: "function",
				Function: &llms.FunctionDefinition{
					Name:        t.Name(),
					Description: t.Description(),
					Parameters:  t.Parameters(),
					Strict:      false,
				},
			})
		}
	}
	return res
}

type OpenaiTool interface {
	tools.Tool
	Parameters() map[string]any
}
