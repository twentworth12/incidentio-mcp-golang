package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// IncidentRole represents an incident role
type IncidentRole struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Shortform    string `json:"shortform"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
	RoleType     string `json:"role_type"`
	Required     bool   `json:"required"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// User represents a user in incident.io (expanded from existing definition)
type UserDetailed struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	SlackUserID   string     `json:"slack_user_id,omitempty"`
	Role          string     `json:"role"`
	BaseRole      *BaseRole  `json:"base_role,omitempty"`
	CustomRoles   []Role     `json:"custom_roles,omitempty"`
}

// BaseRole represents a base role
type BaseRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// Role represents a custom role
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

// ListIncidentRolesOptions represents options for listing incident roles
type ListIncidentRolesOptions struct {
	PageSize int
	After    string
}

// ListIncidentRolesResponse represents the response from listing incident roles
type ListIncidentRolesResponse struct {
	IncidentRoles []IncidentRole `json:"incident_roles"`
	ListResponse
}

// ListUsersOptions represents options for listing users
type ListUsersOptions struct {
	PageSize int
	After    string
	Email    string // Filter by email
}

// ListUsersResponse represents the response from listing users
type ListUsersResponse struct {
	Users []UserDetailed `json:"users"`
	ListResponse
}

// ListIncidentRoles retrieves a list of incident roles
func (c *Client) ListIncidentRoles(opts *ListIncidentRolesOptions) (*ListIncidentRolesResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
	}

	respBody, err := c.doRequest("GET", "/incident_roles", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListIncidentRolesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListUsers retrieves a list of users (with automatic pagination)
func (c *Client) ListUsers(opts *ListUsersOptions) (*ListUsersResponse, error) {
	allUsers := []UserDetailed{}
	pageSize := 250 // Use max page size
	after := ""
	
	// If email filter is provided, don't paginate (API filters server-side)
	if opts != nil && opts.Email != "" {
		params := url.Values{}
		params.Set("page_size", strconv.Itoa(pageSize))
		params.Set("email", opts.Email)
		
		respBody, err := c.doRequest("GET", "/users", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListUsersResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return &response, nil
	}
	
	// For non-filtered requests, paginate through all users
	maxPages := 10 // Safety limit to prevent infinite loops
	for page := 0; page < maxPages; page++ {
		params := url.Values{}
		params.Set("page_size", strconv.Itoa(pageSize))
		if after != "" {
			params.Set("after", after)
		}

		respBody, err := c.doRequest("GET", "/users", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListUsersResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allUsers = append(allUsers, response.Users...)
		
		// Check if there are more pages
		if response.PaginationMeta.After == "" || len(response.Users) == 0 {
			break
		}
		after = response.PaginationMeta.After
	}
	
	// Return combined results
	return &ListUsersResponse{
		Users: allUsers,
		ListResponse: ListResponse{
			PaginationMeta: struct {
				After      string `json:"after,omitempty"`
				PageSize   int    `json:"page_size"`
				TotalCount int    `json:"total_count"`
			}{
				PageSize:   pageSize,
				TotalCount: 0, // API doesn't provide total count
			},
		},
	}, nil
}