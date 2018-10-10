package config

import (
	"encoding/json"
	"strings"

	v "github.com/go-ozzo/ozzo-validation"
)

// Tenant struct
type Tenant struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

// Audit - Audit message for entity
func (tenant *Tenant) Audit() string {
	data, _ := json.Marshal(tenant)
	return string(data)
}

// Validate - Validate fields
func (tenant *Tenant) Validate() error {
	var fieldRules []*v.FieldRules
	// trim space
	tenant.Name = strings.TrimSpace(tenant.Name)
	fieldRules = append(fieldRules, v.Field(&tenant.Name, v.Required.Error("key_name_required"), v.Length(1, 255).Error("key_name_length")))
	fieldRules = append(fieldRules, v.Field(&tenant.Description, v.Length(0, 255).Error("key_description_length")))
	return v.ValidateStruct(tenant, fieldRules...)
}

//SetData - Id, tenant id and user name
func (tenant *Tenant) SetData(id string, tenantID string, userName string) {}
