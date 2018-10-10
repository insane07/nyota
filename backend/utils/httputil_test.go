package utils

import (
	"nyota/backend/i18n"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"bytes"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrepareMapWithTenantInfo(t *testing.T) {
	s := &model.SessionContext{
		User: &model.UserContext{
			UserName: "admin",
			TenantId: "T1",
		},
		Err: nil,
	}
	data := PrepareMapWithTenantInfo(s)

	if len(data) != 1 {
		t.Errorf("Expected map with 1 element...")
	}
	if data["tenant_id"] != s.User.TenantId {
		t.Errorf("Tenant Id not matching...")
	}
}

func TestInvalidInputError(t *testing.T) {

	var jsonStr = []byte(`{"name":"segment1", "aug_methods":[{},{}]}`)
	req := httptest.NewRequest("POST", "/api/v1/segment", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	u := &model.UserContext{
		UserName: "admin",
		TenantId: "T1",
	}
	session := model.SessionContext{User: u, Lang: "en-US", Err: nil}
	session.TFunc = i18n.Translate(&session)

	//var ts config.Segment
	DecodeAndValidate(&session, w, req, nil)

	if session.Err.Code != 500 {
		t.Errorf("Expected error code: 500, but actual error code is :%d", session.Err.Code)
	}
}

func TestInvalidInput(t *testing.T) {

	var jsonStr = []byte(`{"name":"segment1", "aug_methods":[{},{}]}`)
	req := httptest.NewRequest("PUT", "/api/v1/segment", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	u := &model.UserContext{
		UserName: "admin",
		TenantId: "T1",
	}
	session := model.SessionContext{User: u, Lang: "en-US", Err: nil}
	session.TFunc = i18n.Translate(&session)

	var ts config.Segment
	DecodeAndValidate(&session, w, req, &ts)

	if session.Err.Code != 422 {
		t.Errorf("Expected error code: 422, but actual error code is :%d", session.Err.Code)
	}
	if !strings.Contains(session.Err.Message, "Collector Id must be specified") {
		t.Errorf("Expected message : 'Collector Id must be specified', but actual error message is :%s", session.Err.Message)
	}
}

func TestInValidInput2(t *testing.T) {

	var jsonStr = []byte(`{"id":5,"tenant_id":"t1","collector_id":"global","name":"aaaaa","description":"",
		"subnets":["10.2.51.0/24"],"aug_methods":[{"description":"","added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:18:20.452514+05:30"}],
		"added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:10:44.04045+05:30","added_by":"","updated_by":""}`)
	req := httptest.NewRequest("POST", "/api/v1/segment", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	u := &model.UserContext{
		UserName: "admin",
		TenantId: "T1",
	}
	session := model.SessionContext{User: u, Lang: "en-US", Err: nil}
	session.TFunc = i18n.Translate(&session)

	var ts config.Segment
	DecodeAndValidate(&session, w, req, &ts)

	if session.Err == nil {
		t.Errorf("There must not be any errors, but actual error code is :%d", session.Err.Code)
	}
	if session.Err.Code != 422 {
		t.Errorf("Expected error code: 422, but actual error code is :%d", session.Err.Code)
	}
}

func TestValidInput(t *testing.T) {

	var jsonStr = []byte(`{"id":5,"tenant_id":"t1","collector_id":"global","name":"aaaaa","description":"",
		"subnets":["10.2.51.0/24"],"aug_methods":[{"id":1,"tenant_id":"T1","name":"test-nmap","description":"",
		"type":"NMAP","config":{},"added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:18:20.452514+05:30"}],
		"added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:10:44.04045+05:30","added_by":"","updated_by":""}`)
	req := httptest.NewRequest("POST", "/api/v1/segment", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	u := &model.UserContext{
		UserName: "admin",
		TenantId: "T1",
	}
	session := model.SessionContext{User: u, Lang: "en-US", Err: nil}
	session.TFunc = i18n.Translate(&session)

	var ts config.Segment
	DecodeAndValidate(&session, w, req, &ts)
	if session.Err != nil {
		t.Errorf("There must not be any errors, but actual error code is :%d", session.Err.Code)
	}
}

func TestAuditData(t *testing.T) {

	var jsonStr = []byte(`{"id":5,"tenant_id":"t1","collector_id":"global","name":"aaaaa","description":"",
		"subnets":["10.2.51.0/24"],"aug_methods":[{"id":1,"tenant_id":"T1","name":"test-nmap","description":"",
		"type":"NMAP","config":{},"added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:18:20.452514+05:30"}],
		"added_at":"0001-01-01T05:53:28+05:53","updated_at":"2018-02-19T12:10:44.04045+05:30","added_by":"","updated_by":""}`)
	req := httptest.NewRequest("POST", "/api/v1/segment", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	u := &model.UserContext{
		UserName: "admin",
		TenantId: "T1",
	}
	session := model.SessionContext{User: u, Lang: "en-US", Err: nil}
	session.TFunc = i18n.Translate(&session)

	var ts config.Segment
	DecodeAndValidate(&session, w, req, &ts)

	if session.Err != nil {
		t.Errorf("There must not be any errors, but actual error code is :%d", session.Err.Code)
	}
	if session.AuditData == "" {
		t.Errorf("There must be an Audit data for successful operation")
	}
	if !strings.Contains(session.AuditData, "aaaaa") {
		t.Errorf("Invalid Audit data. aaaaa should be present. Found = %s", session.AuditData)
	}
	if !strings.Contains(session.AuditData, "10.2.51.0/24") {
		t.Errorf("Invalid Audit data. 10.2.51.0/24 subnet should be present. Found = %s", session.AuditData)
	}
}

func TestAddValidationError(t *testing.T) {
	s := &model.SessionContext{
		User: &model.UserContext{
			UserName: "admin",
			TenantId: "T1",
		},
		Err: nil,
	}
	err := errors.New("Error")
	addValidationErrors(s, err)

	if s.Err == nil {
		t.Errorf("Expected error...")
	}
	if s.Err.Message != "Error" {
		t.Errorf("Error not matching...")
	}
}

// SetNotFoundError - Sets error to session and handled generically.
func TestSetNotFoundError(t *testing.T) {

	s := &model.SessionContext{
		User: &model.UserContext{
			UserName: "admin",
			TenantId: "T1",
		},
		Err: nil,
	}
	SetNotFoundError(s)
	if s.Err == nil {
		t.Errorf("Expected error...")
	}
	if s.Err.Type != notFoundError {
		t.Errorf("Error not found...")
	}
}

func TestSomethingWentWrong(t *testing.T) {

	s := &model.SessionContext{
		User: &model.UserContext{
			UserName: "admin",
			TenantId: "T1",
		},
		Err: nil,
	}
	SetSomethingWrong(s)
	if s.Err == nil {
		t.Errorf("Expected error...")
	}
	if s.Err.Type != unkownError {
		t.Errorf("Error not found...")
	}
}

func TestParsingError(t *testing.T) {

	s := &model.SessionContext{
		User: &model.UserContext{
			UserName: "admin",
			TenantId: "T1",
		},
		Err: nil,
	}
	s.TFunc = i18n.Translate(s)
	m := make(map[string]interface{})
	m["abcd"] = []interface{}{"abc", "xyz"}
	parseMap(s, m)

	m["abc"] = 1
	parseMap(s, m)

	n := []interface{}{m, m}
	parseArray(s, n)

	n = []interface{}{n, n}
	parseArray(s, n)

	n = []interface{}{1, 2}
	parseArray(s, n)
}
