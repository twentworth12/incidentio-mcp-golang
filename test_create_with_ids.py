#!/usr/bin/env python3
import json
import subprocess
import time

commands = [
    {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}},
    {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {
        "name": "create_incident", 
        "arguments": {
            "name": f"Test Incident {int(time.time())}",
            "summary": "Test incident created via MCP server",
            "severity_id": "01JAR1BCBHF4ZDPZ0HQKJEVZHW",  # Major
            "incident_type_id": "01JAR1BCBH5NN4D1SPFXPTVJGM",  # Default
            "mode": "standard",
            "visibility": "public"
        }
    }}
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
            # Also print stderr for more details
            stderr = proc.stderr.read()
            if stderr:
                print("\nDEBUG INFO:")
                print(stderr)
        elif 'result' in resp:
            print("SUCCESS! Incident created:")
            incident = json.loads(resp['result']['content'][0]['text'])
            print(f"- ID: {incident.get('id')}")
            print(f"- Reference: {incident.get('reference')}")
            print(f"- Name: {incident.get('name')}")
            print(f"- Status: {incident.get('incident_status', {}).get('name')}")
            print(f"- URL: https://app.incident.io/incidents/{incident.get('id')}")