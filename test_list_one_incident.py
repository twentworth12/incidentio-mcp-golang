#!/usr/bin/env python3
import json
import subprocess

# List one incident to see its structure
commands = [
    {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}},
    {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "list_incidents", "arguments": {"page_size": 1}}}
]

proc = subprocess.Popen(['./start-mcp-server.sh'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True)
for cmd in commands:
    proc.stdin.write(json.dumps(cmd) + '\n')
proc.stdin.close()

for line in proc.stdout:
    resp = json.loads(line)
    if resp.get('id') == 2 and 'result' in resp:
        incidents = json.loads(resp['result']['content'][0]['text'])
        if incidents.get('incidents'):
            inc = incidents['incidents'][0]
            print("Sample incident structure:")
            print(f"- Name: {inc.get('name')}")
            print(f"- Mode: {inc.get('mode')}")
            print(f"- Visibility: {inc.get('visibility')}")
            print(f"- Severity ID: {inc.get('severity', {}).get('id')}")
            print(f"- Incident Type ID: {inc.get('incident_type', {}).get('id')}")
            print(f"- Status ID: {inc.get('incident_status', {}).get('id')}")
            
            print("\nFull severity object:", json.dumps(inc.get('severity'), indent=2))
            print("\nFull incident type object:", json.dumps(inc.get('incident_type'), indent=2))