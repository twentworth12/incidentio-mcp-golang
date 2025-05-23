# Testing Guide

This project includes comprehensive tests for all incident.io API endpoints.

## Running Tests

### Unit Tests

Run all unit tests:
```bash
make test-unit
```

Or run Go tests directly:
```bash
go test -v ./internal/incidentio/...
```

### API Integration Tests

Test actual API endpoints (requires `INCIDENT_IO_API_KEY`):
```bash
make test-api
```

Or run the test script directly:
```bash
./test_endpoints.sh
```

### Full Integration Test

Run unit tests, build the server, and test API endpoints:
```bash
make test-integration
```

## Test Coverage

### Unit Tests

The following endpoints have comprehensive unit tests:

#### Workflows (`workflows_test.go`)
- ✅ ListWorkflows - with pagination, empty results, error handling
- ✅ GetWorkflow - success and not found cases
- ✅ UpdateWorkflow - name, enabled status, and state updates

#### Alert Routes (`alert_routes_test.go`)
- ✅ ListAlertRoutes - with pagination and empty results
- ✅ GetAlertRoute - success and not found cases
- ✅ CreateAlertRoute - with conditions, escalations, and templates
- ✅ UpdateAlertRoute - all field updates

#### Alert Sources (`alert_sources_test.go`)
- ✅ ListAlertSources - with pagination and filtering

#### Alert Events (`alert_events_test.go`)
- ✅ CreateAlertEvent - basic, with deduplication, metadata, and resolution

#### Incidents (`incidents_test.go`)
- ✅ ListIncidents - with filters for status and severity
- ✅ GetIncident - with role assignments
- ✅ CreateIncident - with severity, status, and role assignments
- ✅ UpdateIncident - name, summary, status, and severity

### Integration Tests

The `test_endpoints.sh` script tests the following live endpoints:
- Incidents (list, get)
- Incident statuses
- Workflows (list, get)
- Alerts (list, get)
- Alert sources and routes
- Actions (list, get)
- Users and roles

## Test Structure

### Mock HTTP Client

Tests use a mock HTTP client that:
- Validates request methods and paths
- Checks authorization headers
- Verifies query parameters
- Returns predefined responses

### Test Helpers

Common test utilities in `client_test.go`:
- `NewTestClient()` - Creates a client with mock transport
- `mockResponse()` - Creates mock HTTP responses
- `assertNoError()` - Fails test if error occurs
- `assertError()` - Fails test if no error occurs
- `assertEqual()` - Compares expected and actual strings

## Adding New Tests

When adding new endpoints:

1. Create a test file: `internal/incidentio/{endpoint}_test.go`
2. Test all CRUD operations with:
   - Success cases
   - Error cases (404, 500, etc.)
   - Edge cases (empty results, pagination)
3. Verify request formatting:
   - HTTP method
   - URL path
   - Query parameters
   - Request body
4. Add integration test to `test_endpoints.sh`

## Example Test

```go
func TestListResources(t *testing.T) {
    tests := []struct {
        name           string
        params         *ListResourcesParams
        mockResponse   string
        mockStatusCode int
        wantError      bool
        expectedCount  int
    }{
        {
            name:           "successful list",
            params:         &ListResourcesParams{PageSize: 10},
            mockResponse:   `{"resources": [{"id": "123"}]}`,
            mockStatusCode: http.StatusOK,
            wantError:      false,
            expectedCount:  1,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := &MockHTTPClient{
                DoFunc: func(req *http.Request) (*http.Response, error) {
                    // Verify request
                    assertEqual(t, "GET", req.Method)
                    return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
                },
            }
            
            client := NewTestClient(mockClient)
            result, err := client.ListResources(tt.params)
            
            if tt.wantError {
                assertError(t, err)
                return
            }
            
            assertNoError(t, err)
            // Assert results
        })
    }
}
```

## Troubleshooting

### Test Failures

1. **Path mismatches**: Check if the API uses `/v2` prefix
2. **Parameter names**: Verify field names match API documentation
3. **HTTP methods**: Some updates use POST to `/actions/edit` endpoints
4. **Time imports**: Remove unused time imports in test files

### Running Specific Tests

Run a single test:
```bash
go test -v -run TestListWorkflows ./internal/incidentio/
```

Run tests for one package:
```bash
go test -v ./internal/incidentio/
```

## CI/CD Integration

Add to your CI pipeline:
```yaml
- name: Run tests
  run: |
    go test -v ./...
    make test-integration
```