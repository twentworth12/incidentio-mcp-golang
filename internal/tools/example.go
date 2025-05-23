package tools

import "fmt"

type ExampleTool struct{}

func (t *ExampleTool) Name() string {
	return "example_tool"
}

func (t *ExampleTool) Description() string {
	return "An example tool that echoes the input message"
}

func (t *ExampleTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The message to echo",
			},
		},
		"required": []string{"message"},
	}
}

func (t *ExampleTool) Execute(args map[string]interface{}) (string, error) {
	message, ok := args["message"].(string)
	if !ok {
		return "", fmt.Errorf("message parameter is required")
	}
	return fmt.Sprintf("Echo: %s", message), nil
}