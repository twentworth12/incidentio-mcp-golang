package incidentio

import (
	"net/http"
	"testing"
)

func TestListWorkflows(t *testing.T) {
	tests := []struct {
		name           string
		params         *ListWorkflowsParams
		mockResponse   string
		mockStatusCode int
		wantError      bool
		expectedCount  int
	}{
		{
			name:   "successful list workflows",
			params: &ListWorkflowsParams{PageSize: 10},
			mockResponse: `{
				"workflows": [
					{
						"id": "wf_123",
						"name": "Test Workflow",
						"trigger": "incident.created",
						"enabled": true,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					}
				],
				"pagination_info": {
					"page_size": 10
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
			expectedCount:  1,
		},
		{
			name:           "empty workflows list",
			params:         nil,
			mockResponse:   `{"workflows": [], "pagination_info": {"page_size": 25}}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
			expectedCount:  0,
		},
		{
			name:           "API error",
			params:         &ListWorkflowsParams{PageSize: 10},
			mockResponse:   `{"error": "Internal server error"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantError:      true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Verify request
					assertEqual(t, "GET", req.Method)
					assertEqual(t, "Bearer test-api-key", req.Header.Get("Authorization"))
					
					// Check query params if provided
					if tt.params != nil && tt.params.PageSize > 0 {
						assertEqual(t, "10", req.URL.Query().Get("page_size"))
					}
					
					return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
				},
			}
			
			client := NewTestClient(mockClient)
			result, err := client.ListWorkflows(tt.params)
			
			if tt.wantError {
				assertError(t, err)
				return
			}
			
			assertNoError(t, err)
			if len(result.Workflows) != tt.expectedCount {
				t.Errorf("expected %d workflows, got %d", tt.expectedCount, len(result.Workflows))
			}
			
			if tt.expectedCount > 0 {
				assertEqual(t, "wf_123", result.Workflows[0].ID)
				assertEqual(t, "Test Workflow", result.Workflows[0].Name)
			}
		})
	}
}

func TestGetWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		workflowID     string
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:       "successful get workflow",
			workflowID: "wf_123",
			mockResponse: `{
				"workflow": {
					"id": "wf_123",
					"name": "Test Workflow",
					"trigger": "incident.created",
					"enabled": true,
					"state": {"key": "value"},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z"
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
		{
			name:           "workflow not found",
			workflowID:     "wf_nonexistent",
			mockResponse:   `{"error": "Workflow not found"}`,
			mockStatusCode: http.StatusNotFound,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					assertEqual(t, "GET", req.Method)
					assertEqual(t, "/workflows/"+tt.workflowID, req.URL.Path)
					return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
				},
			}
			
			client := NewTestClient(mockClient)
			workflow, err := client.GetWorkflow(tt.workflowID)
			
			if tt.wantError {
				assertError(t, err)
				return
			}
			
			assertNoError(t, err)
			assertEqual(t, tt.workflowID, workflow.ID)
			assertEqual(t, "Test Workflow", workflow.Name)
			if workflow.State["key"] != "value" {
				t.Error("expected state to contain key:value")
			}
		})
	}
}

func TestUpdateWorkflow(t *testing.T) {
	tests := []struct {
		name           string
		workflowID     string
		request        *UpdateWorkflowRequest
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:       "successful update workflow",
			workflowID: "wf_123",
			request: &UpdateWorkflowRequest{
				Name:    "Updated Workflow",
				Enabled: boolPtr(false),
			},
			mockResponse: `{
				"workflow": {
					"id": "wf_123",
					"name": "Updated Workflow",
					"trigger": "incident.created",
					"enabled": false,
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-02T00:00:00Z"
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
		{
			name:       "update workflow state",
			workflowID: "wf_123",
			request: &UpdateWorkflowRequest{
				State: map[string]interface{}{
					"new_key": "new_value",
				},
			},
			mockResponse: `{
				"workflow": {
					"id": "wf_123",
					"name": "Test Workflow",
					"trigger": "incident.created",
					"enabled": true,
					"state": {"new_key": "new_value"},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-02T00:00:00Z"
				}
			}`,
			mockStatusCode: http.StatusOK,
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					assertEqual(t, "PATCH", req.Method)
					assertEqual(t, "/workflows/"+tt.workflowID, req.URL.Path)
					return mockResponse(tt.mockStatusCode, tt.mockResponse), nil
				},
			}
			
			client := NewTestClient(mockClient)
			workflow, err := client.UpdateWorkflow(tt.workflowID, tt.request)
			
			if tt.wantError {
				assertError(t, err)
				return
			}
			
			assertNoError(t, err)
			assertEqual(t, tt.workflowID, workflow.ID)
			
			// Verify updates were applied
			if tt.request.Name != "" {
				assertEqual(t, tt.request.Name, workflow.Name)
			}
			if tt.request.Enabled != nil {
				if workflow.Enabled != *tt.request.Enabled {
					t.Errorf("expected enabled to be %v, got %v", *tt.request.Enabled, workflow.Enabled)
				}
			}
		})
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}