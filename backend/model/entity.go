package model

import (
	"encoding/json"
)

type Event struct {
	UUID        string    `json:"uuid"`
	Data        EventData `json:"data"`
	CPPMVersion string    `json:"cppm_version"`
	TenantID    string    `json:"tenant_id"`
}

func (event Event) MarshalBinary() ([]byte, error) {
	return json.Marshal(event)
}

type EventData struct {
	EntityName string      `json:"Entity_name"`
	CccID      int         `json:"ccc_id"`
	CppmID     int         `json:"cppm_id"`
	URI        string      `json:"uri"`
	Method     string      `json:"method"`
	Payload    interface{} `json:"payload"`
}

func (eventData EventData) MarshalBinary() ([]byte, error) {
	return json.Marshal(eventData)
}

type SimpleEditDataStruct struct {
	Config []DynamicUIField  `json:"uiconfig"`
	Data   interface{}       `json:"data"`
	TOC    UIFormExtraFields `json:"toc,omitempty"`
	Audit  UIFormExtraFields `json:"audit,omitempty"`
	Action UIFormExtraFields `json:"action,omitempty"`
}

// DynamicUIField - Form fields structure
type DynamicUIField struct {
	ControlType string            `json:"control_type"`
	Type        string            `json:"type"`
	Key         string            `json:"key"`
	Label       string            `json:"label"`
	Value       interface{}       `json:"value"`
	Required    bool              `json:"required"`
	Order       int               `json:"order"`
	IsDisabled  bool              `json:"is_disabled"`
	Options     []FormOptions     `json:"options,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Hide        bool              `json:"hide"`
	MinLength   int               `json:"min_length"`
	MaxLength   int               `json:"max_length"`
}

type UIFormExtraFields struct {
	Header string      `json:"header"`
	Show   bool        `json:"show"`
	Fields interface{} `json:"fields,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// FormOptions - Form list option structure
type FormOptions struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UIFormExtraFieldsDetails struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Icon        string      `json:"icon,omitempty"`
	Target      string      `json:"target,omitempty"`
	Show        bool        `json:"show"`
	Value       interface{} `json:"value,omitempty"`
	Menu        interface{} `json:"menu,omitempty"`
}

type VersionEntity interface {
	Convert() VersionEntity // convert object from one version to previous
	GetVersion() int        // get last updated version
}

// DataInfo - Data option
type DataInfo struct {
	Data string `json:"data"`
}
