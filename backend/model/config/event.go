package config

import (
	"encoding/json"
	"strconv"
	"time"
)

// Role - CPPM Role
type Event struct {
	ID            int                    `db:"id" json:"event_id"`
	Name          string                 `db:"name" json:"event_name"`
	Description   string                 `db:"description" json:"event_description"`
	EventDate     time.Time              `db:"event_date" json:"event_date"`
	Detail        map[string]interface{} `db:"detail" json:"event_destail`
	TenantID      string                 `db:"tenant_id" json:"tenant_id"`
	AddedAt       time.Time              `db:"added_at" json:"added_at"`
	UpdatedAt     time.Time              `db:"updated_at" json:"updated_at"`
	AddedAtEpoc   int64                  `db:"-" json:"added_at_epoc"`
	UpdatedAtEpoc int64                  `db:"-" json:"updated_at_epoc"`
}

type EventList struct {
	Events []*Event `json:"events"`
}

// UserEvent - CPPM Role vs Cluster details
type UserEvent struct {
	TenantID  string `db:"tenant_id" json:"tenant_id"`
	ClusterID int    `db:"cluster_id" json:"cluster_id"`
	RoleID    int    `db:"role_id" json:"role_id"`
	CppmID    int    `db:"cppm_id" json:"cppm_id"`
}

// Audit - Audit message for entity
func (event *Event) Audit() string {
	data, _ := json.Marshal(event)
	return string(data)
}

// Validate - Validate fields
func (event *Event) Validate() error {
	return nil
}

//SetData - Id, Cluster id and user name
func (event *Event) SetData(id string, tenantID string, userName string) {
	event.ID, _ = strconv.Atoi(id)
	event.TenantID = tenantID
}
