package api

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/utils"
	"goprizm/httputils"
	"net/http"

	"github.com/gorilla/mux"
)

func (svc *Service) getTenants(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All Tenants invoked...")
	data, err := svc.Store.GetTenants(s)
	if err != nil {
		logutil.Errorf(s, "Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		httputils.ServeJSON(w, data)
	}
}

func (svc *Service) getTenantById(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get Tenant By Id... Id=%v", id)
	data, err := svc.Store.GetTenantById(s, id)
	if err != nil {
		logutil.Errorf(s, "Get Tenant Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		httputils.ServeJSON(w, data)
	}
}

func (svc *Service) UpsertTenant(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Add / Update Tenant Invoked")
	var tenant config.Tenant
	utils.DecodeAndValidate(s, w, req, &tenant)
	if nil != s.Err {
		return
	}
	logutil.Debugf(s, "Tenant object - %v ", tenant)
	_, err := svc.Store.UpsertTenant(s, &tenant)
	if err != nil {
		logutil.Errorf(s, "Upsert Tenant Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (svc *Service) DeleteTenant(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Delete Tenant by id... Id=", id)
	err := svc.Store.DeleteTenantById(s, id)
	if err != nil {
		logutil.Errorf(s, "Delete Tenant Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
