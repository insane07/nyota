package uicomponent

import (
	"nyota/backend/model"
	"nyota/backend/model/config"
)

func GetRoleConfigFormatter(s *model.SessionContext, role *config.Role) *model.SimpleEditDataStruct {
	ds := model.SimpleEditDataStruct{}
	ds.Config = getRoleConfigUIControl(s, role)
	ds.Audit = getRoleConfigAudit(s, role)
	ds.Action = getRoleAction(s)
	ds.Data = role
	return &ds
}

func getRoleConfigAudit(s *model.SessionContext, role *config.Role) model.UIFormExtraFields {
	var auditArr []interface{}
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_at"), Value: role.AddedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_created_by"), Value: s.User.UserName, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_at"), Value: role.UpdatedAtEpoc, Show: true})
	auditArr = append(auditArr, model.UIFormExtraFieldsDetails{Name: s.TFunc("key_updated_by"), Value: s.User.UserName, Show: true})

	audit := model.UIFormExtraFields{Header: s.TFunc("key_add_edit_changes_saved"), Show: true}
	audit.Fields = auditArr
	return audit
}

func getRoleConfigUIControl(s *model.SessionContext, role *config.Role) []model.DynamicUIField {
	var configFormatterList []model.DynamicUIField
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "name", Label: s.TFunc("name"),
		Order: 1, Required: true, Value: role.Name, MinLength: 0, MaxLength: 100})
	configFormatterList = append(configFormatterList, model.DynamicUIField{
		ControlType: "textbox", Key: "description", Label: s.TFunc("description"),
		Order: 2, Required: true, Value: role.Description, MinLength: 0, MaxLength: 100})
	return configFormatterList
}

func GetRoleGridViewColumns(s *model.SessionContext) []config.RoleGridColumn {
	return []config.RoleGridColumn{
		config.RoleGridColumn{Label: s.TFunc("key_name"), Field: "name", CanSort: false, CanFilter: true, Width: 0},
		config.RoleGridColumn{Label: s.TFunc("description"), Field: "description", CanSort: false, CanFilter: true, Width: 0}}
}

func getRoleAction(s *model.SessionContext) model.UIFormExtraFields {
	actionMap := make(map[string]interface{})
	actionMap["attach_cluster"] = model.UIFormExtraFieldsDetails{Name: s.TFunc("key_add_edit_cluster"),
		Show: true, Icon: "hpe:add-new"}

	action := model.UIFormExtraFields{Header: s.TFunc("key_add_edit_Tools"), Show: true}
	action.Fields = actionMap
	return action
}
