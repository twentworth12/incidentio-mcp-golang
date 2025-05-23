#!/usr/bin/env python3
import json
import subprocess
import sys

# Test searching for users with different parameters
commands = [
    {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}},
    # Try listing with a larger page size to see more users
    {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "list_users", "arguments": {"page_size": 100}}}
]

proc = subprocess.Popen(['./start-mcp-server.sh'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
for cmd in commands:
    proc.stdin.write(json.dumps(cmd) + '\n')
proc.stdin.close()

found_charlie = False
all_users = []

for line in proc.stdout:
    resp = json.loads(line)
    if resp.get('id') == 2 and 'result' in resp:
        data = json.loads(resp['result']['content'][0]['text'])
        users = data.get('users', [])
        
        print(f"Found {len(users)} users")
        
        for user in users:
            email = user.get('email', '')
            name = user.get('name', '')
            all_users.append(f"{name} ({email})")
            
            if 'charlie' in email.lower() or 'charlie' in name.lower():
                found_charlie = True
                print(f"\nFOUND CHARLIE: {user['name']} - {user['email']} - ID: {user['id']}")
        
        # Also check pagination info
        pagination = data.get('pagination_meta', {})
        if pagination:
            print(f"\nPagination info: {pagination}")

if not found_charlie:
    print("\nCharlie not found in the results.")
    print(f"Total users returned: {len(all_users)}")
    if len(all_users) < 20:  # Show all if not too many
        print("\nAll users found:")
        for u in all_users:
            print(f"  - {u}")

# Print any errors
stderr = proc.stderr.read()
if stderr:
    print(f"\nErrors: {stderr}", file=sys.stderr)