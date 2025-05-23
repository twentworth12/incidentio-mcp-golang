package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListActionsOptions represents options for listing actions
type ListActionsOptions struct {
	PageSize   int
	After      string
	IncidentID string
	Status     []string
}

// ListActionsResponse represents the response from listing actions
type ListActionsResponse struct {
	Actions []Action `json:"actions"`
	ListResponse
}

// ListActions retrieves a list of actions with automatic pagination
func (c *Client) ListActions(opts *ListActionsOptions) (*ListActionsResponse, error) {
	allActions := []Action{}
	pageSize := 250 // Use max page size
	after := ""
	
	// Set up base parameters
	baseParams := url.Values{}
	if opts != nil {
		if opts.IncidentID != "" {
			baseParams.Set("incident_id", opts.IncidentID)
		}
		for _, status := range opts.Status {
			baseParams.Add("status", status)
		}
	}
	
	// Paginate through all results
	maxPages := 10 // Safety limit
	for page := 0; page < maxPages; page++ {
		params := url.Values{}
		// Copy base parameters
		for k, v := range baseParams {
			params[k] = v
		}
		
		params.Set("page_size", strconv.Itoa(pageSize))
		if after != "" {
			params.Set("after", after)
		}

		respBody, err := c.doRequest("GET", "/actions", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListActionsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allActions = append(allActions, response.Actions...)
		
		// Check if there are more pages
		if response.PaginationMeta.After == "" || len(response.Actions) == 0 {
			break
		}
		after = response.PaginationMeta.After
	}
	
	// Return combined results
	return &ListActionsResponse{
		Actions: allActions,
		ListResponse: ListResponse{
			PaginationMeta: struct {
				After      string `json:"after,omitempty"`
				PageSize   int    `json:"page_size"`
				TotalCount int    `json:"total_count"`
			}{
				PageSize:   pageSize,
				TotalCount: 0,
			},
		},
	}, nil
}

// GetAction retrieves a specific action by ID
func (c *Client) GetAction(id string) (*Action, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/actions/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Action Action `json:"action"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Action, nil
}