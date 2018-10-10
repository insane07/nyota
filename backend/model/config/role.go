package config

import (
	"encoding/json"
	"strconv"
	"time"
)

type ExtraParam map[string]interface{}

// Role - CPPM Role
type Role struct {
	ID            int        `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	Description   string     `db:"description" json:"description"`
	TenantID      string     `db:"tenant_id" json:"tenant_id"`
	PermitID      int        `db:"permit_id" json:"permit_id"`
	Clusters      []*Cluster `db:"-" json:"clusters"`
	AddedAt       time.Time  `db:"added_at" json:"added_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	AddedAtEpoc   int64      `db:"-" json:"added_at_epoc"`
	UpdatedAtEpoc int64      `db:"-" json:"updated_at_epoc"`
	Extras        ExtraParam `db:"extras" json:"extras"`
}

type RoleList struct {
	Roles     []*Role          `json:"roles"`
	Structure []RoleGridColumn `json:"structure"`
}

type RoleGridColumn struct {
	Label       string `json:"label"`
	Field       string `json:"value"`
	CanFilter   bool   `json:"can_filter"`
	CanSort     bool   `json:"can_sort"`
	CanAddAsTag bool   `json:"add_as_tag"`
	Width       int    `json:"width"`
}

// RoleCluster - CPPM Role vs Cluster details
type RoleCluster struct {
	TenantID  string `db:"tenant_id" json:"tenant_id"`
	ClusterID int    `db:"cluster_id" json:"cluster_id"`
	RoleID    int    `db:"role_id" json:"role_id"`
	CppmID    int    `db:"cppm_id" json:"cppm_id"`
}

// Audit - Audit message for entity
func (role *Role) Audit() string {
	data, _ := json.Marshal(role)
	return string(data)
}

// Validate - Validate fields
func (role *Role) Validate() error {
	return nil
}

//SetData - Id, Cluster id and user name
func (role *Role) SetData(id string, tenantID string, userName string) {
	role.ID, _ = strconv.Atoi(id)
	role.TenantID = tenantID
}

//URL - CPPM API URL
func (role *Role) URL() string {
	return "https://localhost/tips/api/role"
}

//EntityName - CPPM Entity Name
func (role *Role) EntityName() string {
	return "Role"
}

// Convert - Convert to lower CPPM version Obj
func (role *Role) Convert() *Role66 {
	prevVersionRole := Role66{}
	bytes, _ := json.Marshal(role)
	json.Unmarshal(bytes, &prevVersionRole)
	return &prevVersionRole
}

// GetVersion - Returns Obj for CPPM Version
func (role *Role) GetVersion() int {
	return 670
}

//Role66 Role for CPPM 6.6
type Role66 struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PermitID      int       `json:"permit_id"`
	AddedAt       time.Time `json:"added_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AddedAtEpoc   int64     `json:"added_at_epoc"`
	UpdatedAtEpoc int64     `json:"updated_at_epoc"`
}

// Convert - Convert to lower CPPM version Obj
func (role *Role66) Convert() *Role66 {
	return role
}

// GetVersion - Returns Obj for CPPM Version
func (role *Role66) GetVersion() int {
	return 660
}
