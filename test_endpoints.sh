#!/bin/bash

# Test script for incident.io MCP server endpoints
# This script tests all the endpoints using the actual API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if API key is set
if [ -z "$INCIDENT_IO_API_KEY" ]; then
    echo -e "${RED}Error: INCIDENT_IO_API_KEY environment variable is not set${NC}"
    echo "Please set it with: export INCIDENT_IO_API_KEY=your-api-key"
    exit 1
fi

API_BASE="https://api.incident.io/v2"
HEADERS="Authorization: Bearer $INCIDENT_IO_API_KEY"

echo -e "${YELLOW}Testing incident.io API Endpoints${NC}"
echo "========================================"

# Function to make API request and check response
test_endpoint() {
    local method=$1
    local endpoint=$2
    local description=$3
    local data=$4
    
    echo -n "Testing $description... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "$HEADERS" \
            -H "Content-Type: application/json" \
            "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "$HEADERS" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo -e "${GREEN}✓ OK ($http_code)${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED ($http_code)${NC}"
        echo "Response: $body"
        return 1
    fi
}

# Test Incident Endpoints
echo -e "\n${YELLOW}1. Testing Incident Endpoints${NC}"
echo "------------------------------"

test_endpoint "GET" "/incidents?page_size=5" "List incidents"
test_endpoint "GET" "/incident_statuses" "List incident statuses"

# Get first incident ID for testing (if any exist)
INCIDENT_ID=$(curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ ! -z "$INCIDENT_ID" ]; then
    test_endpoint "GET" "/incidents/$INCIDENT_ID" "Get specific incident"
fi

# Test Workflow Endpoints
echo -e "\n${YELLOW}2. Testing Workflow Endpoints${NC}"
echo "------------------------------"

test_endpoint "GET" "/workflows?page_size=5" "List workflows"

# Get first workflow ID for testing (if any exist)
WORKFLOW_ID=$(curl -s -H "$HEADERS" "$API_BASE/workflows?page_size=1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ ! -z "$WORKFLOW_ID" ]; then
    test_endpoint "GET" "/workflows/$WORKFLOW_ID" "Get specific workflow"
fi

# Test Alert Endpoints
echo -e "\n${YELLOW}3. Testing Alert Endpoints${NC}"
echo "---------------------------"

test_endpoint "GET" "/alerts?page_size=5" "List alerts"
test_endpoint "GET" "/alert_sources?page_size=5" "List alert sources"
test_endpoint "GET" "/alert_routes?page_size=5" "List alert routes"

# Get first alert ID for testing (if any exist)
ALERT_ID=$(curl -s -H "$HEADERS" "$API_BASE/alerts?page_size=1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ ! -z "$ALERT_ID" ]; then
    test_endpoint "GET" "/alerts/$ALERT_ID" "Get specific alert"
fi

# Get first alert route ID for testing (if any exist)
ROUTE_ID=$(curl -s -H "$HEADERS" "$API_BASE/alert_routes?page_size=1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ ! -z "$ROUTE_ID" ]; then
    test_endpoint "GET" "/alert_routes/$ROUTE_ID" "Get specific alert route"
fi

# Test Action Endpoints
echo -e "\n${YELLOW}4. Testing Action Endpoints${NC}"
echo "----------------------------"

test_endpoint "GET" "/actions?page_size=5" "List actions"

# Get first action ID for testing (if any exist)
ACTION_ID=$(curl -s -H "$HEADERS" "$API_BASE/actions?page_size=1" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ ! -z "$ACTION_ID" ]; then
    test_endpoint "GET" "/actions/$ACTION_ID" "Get specific action"
fi

# Test User and Role Endpoints
echo -e "\n${YELLOW}5. Testing User/Role Endpoints${NC}"
echo "-------------------------------"

test_endpoint "GET" "/users?page_size=5" "List users"
test_endpoint "GET" "/incident_roles" "List incident roles"

echo -e "\n${YELLOW}Summary${NC}"
echo "========"
echo -e "${GREEN}All endpoint tests completed!${NC}"
echo ""
echo "Note: This script only tests READ operations."
echo "To test CREATE/UPDATE operations, use the MCP server tools."