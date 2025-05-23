package incidentio

import "time"

// Incident represents an incident in incident.io
type Incident struct {
	ID                      string              `json:"id"`
	Reference               string              `json:"reference"`
	Name                    string              `json:"name"`
	Summary                 string              `json:"summary,omitempty"`
	Permalink               string              `json:"permalink"`
	IncidentStatus          IncidentStatus      `json:"incident_status"`
	Severity                Severity            `json:"severity"`
	IncidentType            IncidentType        `json:"incident_type"`
	Mode                    string              `json:"mode"`
	Visibility              string              `json:"visibility"`
	CreatedAt               time.Time           `json:"created_at"`
	UpdatedAt               time.Time           `json:"updated_at"`
	SlackTeamID             string              `json:"slack_team_id,omitempty"`
	SlackChannelID          string              `json:"slack_channel_id,omitempty"`
	SlackChannelName        string              `json:"slack_channel_name,omitempty"`
	IncidentRoleAssignments []RoleAssignment    `json:"incident_role_assignments"`
	CustomFieldEntries      []CustomFieldEntry  `json:"custom_field_entries"`
	HasDebrief              bool                `json:"has_debrief"`
}

// IncidentStatus represents the status of an incident
type IncidentStatus struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Rank        int       `json:"rank"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Severity represents the severity of an incident
type Severity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rank        int       `json:"rank"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IncidentType represents the type of an incident
type IncidentType struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	IsDefault              bool      `json:"is_default"`
	PrivateIncidentsOnly   bool      `json:"private_incidents_only"`
	CreateInTriage         string    `json:"create_in_triage"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// RoleAssignment represents a role assignment in an incident
type RoleAssignment struct {
	Role struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Shortform    string `json:"shortform"`
		Description  string `json:"description"`
		Instructions string `json:"instructions"`
		RoleType     string `json:"role_type"`
		Required     bool   `json:"required"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"role"`
	Assignee *User `json:"assignee,omitempty"`
}

// User represents a user in incident.io
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CustomFieldEntry represents a custom field entry
type CustomFieldEntry struct {
	CustomField struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		FieldType   string `json:"field_type"`
		Options     []interface{} `json:"options"`
	} `json:"custom_field"`
	Values []interface{} `json:"values"`
}

// Alert represents an alert in incident.io
type Alert struct {
	ID              string            `json:"id"`
	Title           string            `json:"title"`
	Status          string            `json:"status"`
	Source          string            `json:"source"`
	Metadata        map[string]string `json:"metadata"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	MergedIntoAlert *Alert            `json:"merged_into_alert,omitempty"`
	Incident        *Incident         `json:"incident,omitempty"`
}

// Action represents an action in incident.io
type Action struct {
	ID          string     `json:"id"`
	IncidentID  string     `json:"incident_id"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Assignee    *User      `json:"assignee,omitempty"`
}

// Workflow represents a workflow in incident.io
type Workflow struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Trigger     string                 `json:"trigger"`
	Enabled     bool                   `json:"enabled"`
	Runs        []WorkflowRun          `json:"runs,omitempty"`
	State       map[string]interface{} `json:"state,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// WorkflowRun represents a workflow run
type WorkflowRun struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AlertRoute represents an alert route in incident.io
type AlertRoute struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Enabled      bool                   `json:"enabled"`
	Conditions   []AlertCondition       `json:"conditions"`
	Escalations  []EscalationBinding    `json:"escalations"`
	GroupingKeys []string               `json:"grouping_keys,omitempty"`
	Template     map[string]interface{} `json:"template,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// AlertCondition represents a condition for alert routing
type AlertCondition struct {
	Field     string `json:"field"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

// EscalationBinding represents an escalation in alert routing
type EscalationBinding struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

// AlertEvent represents an alert event
type AlertEvent struct {
	ID              string                 `json:"id"`
	AlertSourceID   string                 `json:"alert_source_id"`
	DeduplicationKey string                 `json:"deduplication_key,omitempty"`
	Status          string                 `json:"status"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CreateIncidentRequest represents a request to create an incident
type CreateIncidentRequest struct {
	IdempotencyKey           string                          `json:"idempotency_key"`
	Name                     string                          `json:"name"`
	Summary                  string                          `json:"summary,omitempty"`
	IncidentStatusID         string                          `json:"incident_status_id,omitempty"`
	SeverityID               string                          `json:"severity_id,omitempty"`
	IncidentTypeID           string                          `json:"incident_type_id,omitempty"`
	Mode                     string                          `json:"mode,omitempty"`
	Visibility               string                          `json:"visibility,omitempty"`
	CustomFieldEntries       []CustomFieldEntryRequest       `json:"custom_field_entries,omitempty"`
	IncidentRoleAssignments  []CreateRoleAssignmentRequest   `json:"incident_role_assignments,omitempty"`
	SlackChannelNameOverride string                          `json:"slack_channel_name_override,omitempty"`
}

// CustomFieldEntryRequest represents a custom field entry in create/update requests
type CustomFieldEntryRequest struct {
	CustomFieldID string        `json:"custom_field_id"`
	Values        []interface{} `json:"values"`
}

// CreateRoleAssignmentRequest represents a role assignment in create request
type CreateRoleAssignmentRequest struct {
	IncidentRoleID string `json:"incident_role_id"`
	UserID         string `json:"user_id"`
}

// UpdateIncidentRequest represents a request to update an incident
type UpdateIncidentRequest struct {
	Name                    string                          `json:"name,omitempty"`
	Summary                 string                          `json:"summary,omitempty"`
	IncidentStatusID        string                          `json:"incident_status_id,omitempty"`
	SeverityID              string                          `json:"severity_id,omitempty"`
	CustomFieldEntries      []CustomFieldEntryRequest       `json:"custom_field_entries,omitempty"`
	IncidentRoleAssignments []CreateRoleAssignmentRequest   `json:"incident_role_assignments,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	PaginationMeta struct {
		After      string `json:"after,omitempty"`
		PageSize   int    `json:"page_size"`
		TotalCount int    `json:"total_count"`
	} `json:"pagination_meta"`
}