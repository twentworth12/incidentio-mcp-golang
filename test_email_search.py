#!/usr/bin/env python3
import json
import subprocess

commands = [
    {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}},
    # Test email search
    {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "list_users", "arguments": {"email": "charlie@incident.io"}}}
]

proc = subprocess.Popen(['./start-mcp-server.sh'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
for cmd in commands:
    proc.stdin.write(json.dumps(cmd) + '\n')
proc.stdin.close()

for line in proc.stdout:
    resp = json.loads(line)
    if resp.get('id') == 2:
        if 'error' in resp:
            print("ERROR:", resp['error'])
        elif 'result' in resp:
            print(resp['result']['content'][0]['text'])