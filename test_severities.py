#!/usr/bin/env python3
import json
import subprocess

# Initialize
init_req = {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}

# List severities  
sev_req = {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "list_severities", "arguments": {}}}

# List incident types
types_req = {"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "list_incident_types", "arguments": {}}}

proc = subprocess.Popen(['./start-mcp-server.sh'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True)
proc.stdin.write(json.dumps(init_req) + '\n')
proc.stdin.write(json.dumps(sev_req) + '\n') 
proc.stdin.write(json.dumps(types_req) + '\n')
proc.stdin.close()

for line in proc.stdout:
    resp = json.loads(line)
    if resp.get('id') == 2:
        print("SEVERITIES:")
        print(resp['result']['content'][0]['text'])
    elif resp.get('id') == 3:
        print("\nINCIDENT TYPES:")
        print(resp['result']['content'][0]['text'])