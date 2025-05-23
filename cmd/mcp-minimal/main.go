package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type Message struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue
		}

		var response Message

		switch msg.Method {
		case "initialize":
			response = Message{
				Jsonrpc: "2.0",
				ID:      msg.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{},
					},
					"serverInfo": map[string]interface{}{
						"name":    "incident-mcp",
						"version": "1.0.0",
					},
				},
			}
		case "initialized":
			continue // No response needed
		case "tools/list":
			response = Message{
				Jsonrpc: "2.0",
				ID:      msg.ID,
				Result: map[string]interface{}{
					"tools": []map[string]interface{}{
						{
							"name":        "test_tool",
							"description": "A simple test tool",
							"inputSchema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"message": map[string]interface{}{
										"type":        "string",
										"description": "Test message",
									},
								},
								"required": []string{"message"},
							},
						},
					},
				},
			}
		case "tools/call":
			params := msg.Params.(map[string]interface{})
			response = Message{
				Jsonrpc: "2.0",
				ID:      msg.ID,
				Result: map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": "Tool called with: " + params["name"].(string),
						},
					},
				},
			}
		default:
			response = Message{
				Jsonrpc: "2.0",
				ID:      msg.ID,
				Error: &Error{
					Code:    -32601,
					Message: "Method not found",
				},
			}
		}

		encoder.Encode(response)
	}
}