package uicomponent

import (
	"nyota/backend/model"
	"nyota/backend/model/config"
	"strconv"
)

func GetCPPMNodeConfigFormatter(s *model.SessionContext, cppmNode *config.CppmNode, clusters []*config.Cluster) *model.SimpleEditDataStruct {
	ds := model.SimpleEditDataStruct{}
	ds.Config = getCPPMNodeConfigUIControl(s, cppmNode, clusters)
	ds.Audit = getCppmNodeConfigAudit(s, cppmNode)
	ds.Data = cppmNode
	return &ds
}

func getCPPMNodeConfigUIControl(s *model.SessionContext, cppmNode *config.CppmNode, clusters []*config.Cluster) []model.DynamicUIField {
	var configFormatterList []model.DynamicUIField
	formOptions := []model.FormOptions{}
	colValue := ""

	if nil != clusters && len(clusters) > 0 {
		for _, cluster := range clusters {
			formOptions = append(formOptions, model.FormOptions{Key: strconv.Itoa(cluster.ID), Value: cluster.Name})
		}
		colValue = formOptions[0].Key
	}

	if cppmNode.ClusterID != 0 {
		colValue = strconv.Itoa(cppmNode.ClusterID)
	}

	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "dropdown", Key: "cluster_id", Label: s.TFunc("cluster_id"),
		Order: 1, Required: true, Value: colValue, Options: formOptions})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_standby", Label: s.TFunc("is_standby"),
		Order: 2, Required: true, Value: cppmNode.IsStandBy, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "cppm_version", Label: s.TFunc("cppm_version"),
		Order: 3, Required: true, Value: cppmNode.CppmVersion, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "server_uuid", Label: s.TFunc("server_uuid"),
		Order: 4, Required: true, Value: cppmNode.ServerUUID, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "server_dns_name", Label: s.TFunc("server_dns_name"),
		Order: 5, Required: true, Value: cppmNode.ServerDNSName, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "fqdn", Label: s.TFunc("fqdn"),
		Order: 6, Required: true, Value: cppmNode.Fqdn, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "server_ip", Label: s.TFunc("server_ip"),
		Order: 7, Required: true, Value: cppmNode.ServerIP, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "management_ip", Label: s.TFunc("management_ip"),
		Order: 8, Required: true, Value: cppmNode.ManagementIP, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "ipv6_server_ip", Label: s.TFunc("ipv6_server_ip"),
		Order: 9, Required: false, Value: cppmNode.IPV6ServerIP, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "ipv6_management_ip", Label: s.TFunc("ipv6_management_ip"),
		Order: 10, Required: false, Value: cppmNode.IPV6ManagementIP, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_master", Label: s.TFunc("is_master"),
		Order: 11, Required: true, Value: cppmNode.IsMaster, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "provider_uuid", Label: s.TFunc("provider_uuid"),
		Order: 12, Required: true, Value: cppmNode.ProviderUUID, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "domain_id", Label: s.TFunc("domain_id"),
		Order: 13, Required: false, Value: cppmNode.DomainID, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_profiler_enabled", Label: s.TFunc("is_profiler_enabled"),
		Order: 14, Required: true, Value: cppmNode.IsProfilerEnabled, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_insight_enabled", Label: s.TFunc("is_insight_enabled"),
		Order: 15, Required: true, Value: cppmNode.IsInsightEnabled, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_insight_master", Label: s.TFunc("is_insight_master"),
		Order: 16, Required: true, Value: cppmNode.IsInsightMaster, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_perfmas_enabled", Label: s.TFunc("is_perfmas_enabled"),
		Order: 16, Required: true, Value: cppmNode.IsPerfmasEnabled, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_cloud_tunnel_enabled", Label: s.TFunc("is_cloud_tunnel_enabled"),
		Order: 17, Required: true, Value: cppmNode.IsCloudTunnelEnabled, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "checkbox", Key: "is_ingress_events_enabled", Label: s.TFunc("is_ingress_events_enabled"),
		Order: 18, Required: true, Value: cppmNode.IsIngressEventsEnabled, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "dhcp_span_intf", Label: s.TFunc("dhcp_span_intf"),
		Order: 19, Required: false, Value: cppmNode.DhcpSpanIntf, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "replication_status", Label: s.TFunc("replication_status"),
		Order: 20, Required: false, Value: cppmNode.ReplicationStatus, MinLength: 0, MaxLength: 100})

	return configFormatterList
}

func getCppmNodeConfigAudit(s *model.SessionContext, cppmNode *config.CppmNode) model.UIFormExtraFields {
	var auditArr []interface{}
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_at"), Value: cppmNode.AddedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_by"), Value: s.User.UserName, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_at"), Value: cppmNode.UpdatedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_by"), Value: s.User.UserName, Show: true})

	audit := model.UIFormExtraFields{Header: s.TFunc("key_add_edit_changes_saved"), Show: true}
	audit.Fields = auditArr
	return audit
}
