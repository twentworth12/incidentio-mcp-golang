package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListWorkflowsParams contains optional parameters for listing workflows
type ListWorkflowsParams struct {
	PageSize int
	After    string
}

// ListWorkflowsResponse represents the response from listing workflows
type ListWorkflowsResponse struct {
	Workflows []Workflow `json:"workflows"`
	Pagination struct {
		After  string `json:"after,omitempty"`
		PageSize int `json:"page_size"`
	} `json:"pagination_info"`
}

// ListWorkflows returns all workflows
func (c *Client) ListWorkflows(params *ListWorkflowsParams) (*ListWorkflowsResponse, error) {
	endpoint := "/workflows"
	
	v := url.Values{}
	if params != nil {
		if params.PageSize > 0 {
			v.Set("page_size", fmt.Sprintf("%d", params.PageSize))
		}
		if params.After != "" {
			v.Set("after", params.After)
		}
	}
	
	if len(v) > 0 {
		endpoint = endpoint + "?" + v.Encode()
	}
	
	respBody, err := c.doRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	
	var result ListWorkflowsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &result, nil
}

// GetWorkflow returns a specific workflow by ID
func (c *Client) GetWorkflow(id string) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s", id)
	
	respBody, err := c.doRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Workflow Workflow `json:"workflow"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &result.Workflow, nil
}

// UpdateWorkflowRequest represents a request to update a workflow
type UpdateWorkflowRequest struct {
	Name    string                 `json:"name,omitempty"`
	Enabled *bool                  `json:"enabled,omitempty"`
	State   map[string]interface{} `json:"state,omitempty"`
}

// UpdateWorkflow updates a workflow
func (c *Client) UpdateWorkflow(id string, req *UpdateWorkflowRequest) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s", id)
	
	respBody, err := c.doRequest("PATCH", endpoint, nil, req)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Workflow Workflow `json:"workflow"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &result.Workflow, nil
}