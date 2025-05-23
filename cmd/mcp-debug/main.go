package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func logDebug(msg string, data interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	if data != nil {
		if jsonData, err := json.Marshal(data); err == nil {
			fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, msg, string(jsonData))
		} else {
			fmt.Fprintf(os.Stderr, "[%s] %s: %+v\n", timestamp, msg, data)
		}
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", timestamp, msg)
	}
}

func main() {
	logDebug("MCP Server Starting", nil)
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		logDebug("RECEIVED", line)
		
		if line == "" {
			continue
		}
		
		var request map[string]interface{}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			logDebug("JSON PARSE ERROR", err.Error())
			continue
		}
		
		logDebug("PARSED REQUEST", request)
		
		method, _ := request["method"].(string)
		id := request["id"]
		
		var response map[string]interface{}
		
		switch method {
		case "initialize":
			logDebug("HANDLING INITIALIZE", nil)
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{},
					},
					"serverInfo": map[string]interface{}{
						"name":    "debug-mcp-server",
						"version": "1.0.0",
					},
				},
			}
		case "initialized":
			logDebug("HANDLING INITIALIZED", nil)
			continue
		case "tools/list":
			logDebug("HANDLING TOOLS/LIST", nil)
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"tools": []map[string]interface{}{
						{
							"name":        "debug_tool",
							"description": "Debug test tool",
							"inputSchema": map[string]interface{}{
								"type":       "object",
								"properties": map[string]interface{}{},
							},
						},
					},
				},
			}
		case "tools/call":
			logDebug("HANDLING TOOLS/CALL", request["params"])
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": "Debug tool executed successfully",
						},
					},
				},
			}
		default:
			logDebug("UNKNOWN METHOD", method)
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": "Method not found: " + method,
				},
			}
		}
		
		logDebug("SENDING RESPONSE", response)
		
		output, err := json.Marshal(response)
		if err != nil {
			logDebug("JSON MARSHAL ERROR", err.Error())
			continue
		}
		
		logDebug("RAW OUTPUT", string(output))
		os.Stdout.Write(output)
		os.Stdout.Write([]byte("\n"))
		os.Stdout.Sync()
	}
	
	if err := scanner.Err(); err != nil {
		logDebug("SCANNER ERROR", err.Error())
	}
	
	logDebug("MCP Server Stopping", nil)
}