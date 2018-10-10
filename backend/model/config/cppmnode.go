package config

import (
	"encoding/json"
	"strconv"
	"time"

	v "github.com/go-ozzo/ozzo-validation"

	"nyota/backend/utils"
)

type CppmNode struct {
	ID                       int       `db:"id" json:"id" props:"primary_key=true"`
	TenantID                 string    `db:"tenant_id" json:"tenant_id"`
	ClusterID                int       `db:"cluster_id" json:"cluster_id"`
	IsStandBy                bool      `db:"is_standby" json:"is_standby"`
	CppmVersion              string    `db:"cppm_version" json:"cppm_version"`
	ServerUUID               string    `db:"server_uuid" json:"server_uuid"`
	ServerDNSName            string    `db:"server_dns_name" json:"server_dns_name"`
	Fqdn                     string    `db:"fqdn" json:"fqdn"`
	ServerIP                 string    `db:"server_ip" json:"server_ip"`
	ManagementIP             string    `db:"management_ip" json:"management_ip"`
	IPV6ServerIP             string    `db:"ipv6_server_ip" json:"ipv6_server_ip"`
	IPV6ManagementIP         string    `db:"ipv6_management_ip" json:"ipv6_management_ip"`
	IsMaster                 bool      `db:"is_master" json:"is_master"`
	ProviderUUID             string    `db:"provider_uuid" json:"provider_uuid"`
	DomainID                 int       `db:"domain_id" json:"domain_id"`
	IsProfilerEnabled        bool      `db:"is_profiler_enabled" json:"is_profiler_enabled"`
	IsInsightEnabled         bool      `db:"is_insight_enabled" json:"is_insight_enabled"`
	IsInsightMaster          bool      `db:"is_insight_master" json:"is_insight_master"`
	IsPerfmasEnabled         bool      `db:"is_perfmas_enabled" json:"is_perfmas_enabled"`
	IsCloudTunnelEnabled     bool      `db:"is_cloud_tunnel_enabled" json:"is_cloud_tunnel_enabled"`
	IsIngressEventsEnabled   bool      `db:"is_ingress_events_enabled" json:"is_ingress_events_enabled"`
	DhcpSpanIntf             string    `db:"dhcp_span_intf" json:"dhcp_span_intf"`
	ReplicationStatus        string    `db:"replication_status" json:"replication_status"`
	LastReplicationTimestamp time.Time `db:"last_replication_timestamp" json:"last_replication_timestamp"`
	AddedAt                  time.Time `db:"added_at" json:"added_at"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
	AddedAtEpoc              int64     `db:"-" json:"added_at_epoc"`
	UpdatedAtEpoc            int64     `db:"-" json:"updated_at_epoc"`
}

// Audit message for entity.
func (cppmnode *CppmNode) Audit() string {
	data, _ := json.Marshal(cppmnode)
	return string(data)
}

// Validating fields.
func (cppmnode *CppmNode) Validate() error {
	var fieldRules []*v.FieldRules
	fieldRules = append(fieldRules, v.Field(&cppmnode.CppmVersion, v.Required.Error("cppm_version_is_required")))
	fieldRules = append(fieldRules, v.Field(&cppmnode.ServerIP, v.Required.Error("server_ip_is_invalid"), v.By(utils.ValidateIP)))
	fieldRules = append(fieldRules, v.Field(&cppmnode.ManagementIP, v.Required.Error("mgmt_ip_is_invalid"), v.By(utils.ValidateIP)))
	return v.ValidateStruct(cppmnode, fieldRules...)
}

// Setting ID, Tenant ID and Username.
func (cppmnode *CppmNode) SetData(id string, tenantID string, userName string) {
	cppmnode.ID, _ = strconv.Atoi(id)
	cppmnode.TenantID = tenantID
}
