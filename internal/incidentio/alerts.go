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

// ListAlerts retrieves a list of alerts
func (c *Client) ListAlerts(opts *ListAlertsOptions) (*ListAlertsResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
		for _, status := range opts.Status {
			params.Add("status", status)
		}
	}

	respBody, err := c.doRequest("GET", "/alerts", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListAlertsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
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

// ListAlertsForIncident retrieves alerts for a specific incident
func (c *Client) ListAlertsForIncident(incidentID string, opts *ListAlertsOptions) (*ListAlertsResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
	}

	respBody, err := c.doRequest("GET", fmt.Sprintf("/alerts/incident/%s", incidentID), params, nil)
	if err != nil {
		return nil, err
	}

	var response ListAlertsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}