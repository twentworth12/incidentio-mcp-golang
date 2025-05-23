#!/usr/bin/env python3
"""
Simple test client for the Incident.io MCP server
"""
import json
import subprocess
import sys
import os

def send_mcp_message(process, method, params=None, message_id=1):
    """Send a JSON-RPC message to the MCP server"""
    message = {
        "jsonrpc": "2.0",
        "method": method,
        "id": message_id
    }
    if params:
        message["params"] = params
    
    message_json = json.dumps(message) + "\n"
    process.stdin.write(message_json.encode())
    process.stdin.flush()
    
    # Read response
    response_line = process.stdout.readline()
    if not response_line:
        return None
    
    return json.loads(response_line.decode())

def test_mcp_server():
    """Test the MCP server functionality"""
    # Set environment variable
    env = os.environ.copy()
    env["INCIDENT_IO_API_KEY"] = "inc_4f78bf72eb22ce2f6a58be0feb397a4030e1511c8d0eec234759ef83c00e8690"
    
    # Start the MCP server
    try:
        process = subprocess.Popen(
            ["go", "run", "cmd/mcp-server/main.go"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            env=env,
            text=False
        )
        
        print("🚀 Starting MCP server test...")
        
        # Test 1: Initialize
        print("\n1. Testing initialization...")
        response = send_mcp_message(process, "initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "test-client", "version": "1.0.0"}
        })
        
        if response and "result" in response:
            print("✅ Initialization successful")
            print(f"   Server: {response['result']['serverInfo']['name']} v{response['result']['serverInfo']['version']}")
        else:
            print("❌ Initialization failed")
            return False
        
        # Test 2: List tools
        print("\n2. Testing tools list...")
        response = send_mcp_message(process, "tools/list", {}, 2)
        
        if response and "result" in response:
            tools = response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools:")
            for tool in tools[:5]:  # Show first 5 tools
                print(f"   - {tool['name']}: {tool['description']}")
            if len(tools) > 5:
                print(f"   ... and {len(tools) - 5} more")
        else:
            print("❌ Tools list failed")
            return False
        
        # Test 3: Call a tool (list_incidents)
        print("\n3. Testing list_incidents tool...")
        response = send_mcp_message(process, "tools/call", {
            "name": "list_incidents",
            "arguments": {"page_size": 2}
        }, 3)
        
        if response and "result" in response:
            print("✅ list_incidents tool executed successfully")
            # Parse the JSON response from the tool
            try:
                incidents_data = json.loads(response["result"]["content"][0]["text"])
                incident_count = len(incidents_data.get("incidents", []))
                print(f"   Found {incident_count} incidents")
                if incident_count > 0:
                    first_incident = incidents_data["incidents"][0]
                    print(f"   Example: {first_incident.get('name', 'Unknown')} (Status: {first_incident.get('status', 'Unknown')})")
            except json.JSONDecodeError:
                print("   Raw response received (not JSON)")
        else:
            print("❌ list_incidents tool failed")
            if response and "error" in response:
                print(f"   Error: {response['error']['message']}")
            return False
        
        print("\n🎉 All tests passed! MCP server is working correctly.")
        return True
        
    except FileNotFoundError:
        print("❌ Go not found. Please install Go to run this test.")
        return False
    except Exception as e:
        print(f"❌ Test failed with error: {e}")
        return False
    finally:
        if 'process' in locals():
            process.terminate()
            process.wait()

if __name__ == "__main__":
    success = test_mcp_server()
    sys.exit(0 if success else 1)