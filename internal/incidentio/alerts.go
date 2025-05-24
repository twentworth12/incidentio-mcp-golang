package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListAlertsOptions represents options for listing alerts
type ListAlertsOptions struct {
	PageSize int
	After    string
	Status   []string
}

// ListAlertsResponse represents the response from listing alerts
type ListAlertsResponse struct {
	Alerts []Alert `json:"alerts"`
	ListResponse
}

// ListAlerts retrieves a list of alerts with automatic pagination
func (c *Client) ListAlerts(opts *ListAlertsOptions) (*ListAlertsResponse, error) {
	allAlerts := []Alert{}
	pageSize := 50 // Max page size for alerts is 50
	after := ""
	
	// Set up base parameters
	baseParams := url.Values{}
	if opts != nil {
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

		respBody, err := c.doRequest("GET", "/alerts", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListAlertsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allAlerts = append(allAlerts, response.Alerts...)
		
		// Check if there are more pages
		if response.PaginationMeta.After == "" || len(response.Alerts) == 0 {
			break
		}
		after = response.PaginationMeta.After
	}
	
	// Return combined results
	return &ListAlertsResponse{
		Alerts: allAlerts,
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

// GetAlert retrieves a specific alert by ID
func (c *Client) GetAlert(id string) (*Alert, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/alerts/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Alert Alert `json:"alert"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Alert, nil
}

// ListAlertsForIncident retrieves alerts for a specific incident with automatic pagination
func (c *Client) ListAlertsForIncident(incidentID string, opts *ListAlertsOptions) (*ListAlertsResponse, error) {
	allAlerts := []Alert{}
	pageSize := 50 // Max page size for alerts is 50
	after := ""
	
	// Set up base parameters
	baseParams := url.Values{}
	if opts != nil {
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
		
		params.Set("incident_id", incidentID) // Filter by incident
		params.Set("page_size", strconv.Itoa(pageSize))
		if after != "" {
			params.Set("after", after)
		}

		respBody, err := c.doRequest("GET", "/alerts", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListAlertsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allAlerts = append(allAlerts, response.Alerts...)
		
		// Check if there are more pages
		if response.PaginationMeta.After == "" || len(response.Alerts) == 0 {
			break
		}
		after = response.PaginationMeta.After
	}
	
	// Return combined results
	return &ListAlertsResponse{
		Alerts: allAlerts,
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