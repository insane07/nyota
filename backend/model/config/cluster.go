package config

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
)

// Cluster struct
type Cluster struct {
	ID            int         `db:"id" json:"id" props:"primary_key=true"`
	UUID          string      `db:"uuid" json:"uuid"`
	Name          string      `db:"name" json:"name"`
	Description   string      `db:"description" json:"description"`
	CppmVersion   string      `db:"cppm_version" json:"cppm_version"`
	TenantID      string      `db:"tenant_id" json:"tenant_id"`
	CPPMNodes     []*CppmNode `db:"-" json:"cppm_nodes"`
	AddedAt       time.Time   `db:"added_at" json:"added_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
	AddedAtEpoc   int64       `db:"-" json:"added_at_epoc"`
	UpdatedAtEpoc int64       `db:"-" json:"updated_at_epoc"`
}

// Audit - Audit message for entity
func (cluster *Cluster) Audit() string {
	data, _ := json.Marshal(cluster)
	return string(data)
}

// Validate - Validate fields
func (cluster *Cluster) Validate() error {
	var fieldRules []*v.FieldRules
	// trim space
	cluster.Name = strings.TrimSpace(cluster.Name)
	fieldRules = append(fieldRules, v.Field(&cluster.Name, v.Required.Error("key_name_required"), v.Length(1, 255).Error("key_name_length")))
	fieldRules = append(fieldRules, v.Field(&cluster.Description, v.Length(0, 255).Error("key_description_length")))
	return v.ValidateStruct(cluster, fieldRules...)
}

//SetData - Id, Cluster id and user name
func (cluster *Cluster) SetData(id string, tenantID string, userName string) {
	cluster.ID, _ = strconv.Atoi(id)
	cluster.TenantID = tenantID
}
