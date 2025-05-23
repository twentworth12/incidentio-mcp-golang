package main

import (
	"bufio"
	"encoding/json"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		var request map[string]interface{}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue
		}
		
		method, _ := request["method"].(string)
		id := request["id"]
		
		var response map[string]interface{}
		
		switch method {
		case "initialize":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{},
					},
					"serverInfo": map[string]interface{}{
						"name":    "test-server",
						"version": "1.0.0",
					},
				},
			}
		case "initialized":
			continue
		case "tools/list":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"tools": []map[string]interface{}{
						{
							"name":        "hello",
							"description": "Say hello",
							"inputSchema": map[string]interface{}{
								"type":       "object",
								"properties": map[string]interface{}{},
							},
						},
					},
				},
			}
		case "tools/call":
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": "Hello from Go MCP server!",
						},
					},
				},
			}
		default:
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": "Method not found",
				},
			}
		}
		
		output, _ := json.Marshal(response)
		os.Stdout.Write(output)
		os.Stdout.Write([]byte("\n"))
	}
}