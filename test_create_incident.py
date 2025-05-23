#!/usr/bin/env python3
"""
Test creating an incident through the MCP server
"""
import json
import subprocess
import sys

def send_request(request):
    """Send a request to the MCP server and return the response"""
    proc = subprocess.Popen(
        ['./start-mcp-server.sh'],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    stdout, stderr = proc.communicate(input=json.dumps(request))
    
    if stderr:
        print(f"STDERR: {stderr}", file=sys.stderr)
    
    try:
        return json.loads(stdout)
    except json.JSONDecodeError:
        print(f"Failed to parse response: {stdout}")
        return None

def main():
    print("Testing incident creation through MCP server...")
    
    # Initialize
    print("\n1. Initializing connection...")
    init_request = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {
                "name": "test-client",
                "version": "1.0.0"
            }
        }
    }
    
    response = send_request(init_request)
    print(f"Initialize response: {json.dumps(response, indent=2)}")
    
    # List available tools
    print("\n1.5. Listing available tools...")
    list_tools_request = {
        "jsonrpc": "2.0",
        "id": 15,
        "method": "tools/list",
        "params": {}
    }
    
    response = send_request(list_tools_request)
    if response and 'result' in response:
        print("Available tools:")
        for tool in response['result']['tools']:
            print(f"  - {tool['name']}: {tool['description']}")
    
    # List incident types
    print("\n2. Listing incident types...")
    list_types_request = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "list_incident_types",
            "arguments": {}
        }
    }
    
    response = send_request(list_types_request)
    if response and 'result' in response:
        print("Incident Types:")
        print(response['result']['content'][0]['text'])
    
    # List severities
    print("\n3. Listing severities...")
    list_severities_request = {
        "jsonrpc": "2.0",
        "id": 3,
        "method": "tools/call",
        "params": {
            "name": "list_severities",
            "arguments": {}
        }
    }
    
    response = send_request(list_severities_request)
    if response and 'result' in response:
        print("Severities:")
        print(response['result']['content'][0]['text'])
    
    # List incident statuses
    print("\n4. Listing incident statuses...")
    list_statuses_request = {
        "jsonrpc": "2.0",
        "id": 4,
        "method": "tools/call",
        "params": {
            "name": "list_incident_statuses",
            "arguments": {}
        }
    }
    
    response = send_request(list_statuses_request)
    if response and 'result' in response:
        print("Incident Statuses:")
        print(response['result']['content'][0]['text'])
    
    # First, let's try to create an incident without severity to see what's required
    print("\n5. Creating test incident (attempt 1 - minimal fields)...")
    create_incident_request = {
        "jsonrpc": "2.0",
        "id": 5,
        "method": "tools/call",
        "params": {
            "name": "create_incident",
            "arguments": {
                "name": "Test Incident from MCP",
                "summary": "This is a test incident created through the MCP server",
                "mode": "standard",
                "visibility": "public"
            }
        }
    }
    
    response = send_request(create_incident_request)
    if response:
        if 'error' in response:
            print(f"Error creating incident: {response['error']}")
        elif 'result' in response:
            print("Incident created successfully!")
            print(response['result']['content'][0]['text'])
    
    print("\nTest complete!")

if __name__ == "__main__":
    main()