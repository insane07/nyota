package uicomponent

import (
	"nyota/backend/model"
	"nyota/backend/model/config"
)

func GetClusterConfigFormatter(s *model.SessionContext, cluster *config.Cluster) *model.SimpleEditDataStruct {
	ds := model.SimpleEditDataStruct{}
	ds.Config = getClusterConfigUIControl(s, cluster)
	ds.Audit = getClusterConfigAudit(s, cluster)
	ds.Data = cluster
	return &ds
}

func getClusterConfigAudit(s *model.SessionContext, cluster *config.Cluster) model.UIFormExtraFields {
	var auditArr []interface{}
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_at"), Value: cluster.AddedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_by"), Value: s.User.UserName, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_at"), Value: cluster.UpdatedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_by"), Value: s.User.UserName, Show: true})

	audit := model.UIFormExtraFields{Header: s.TFunc("key_add_edit_changes_saved"), Show: true}
	audit.Fields = auditArr
	return audit
}

func getClusterConfigUIControl(s *model.SessionContext, cluster *config.Cluster) []model.DynamicUIField {
	var configFormatterList []model.DynamicUIField
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "name", Label: s.TFunc("name"),
		Order: 1, Required: true, Value: cluster.Name, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "description", Label: s.TFunc("description"),
		Order: 2, Required: true, Value: cluster.Description, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "cppm_version", Label: s.TFunc("cppm_version"),
		Order: 3, Required: true, Value: cluster.CppmVersion, MinLength: 0, MaxLength: 100})
	return configFormatterList
}
