package api

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/uicomponent"
	"nyota/backend/utils"
	"encoding/json"
	"goprizm/httputils"
	"net/http"

	"github.com/gorilla/mux"
)

func (svc *Service) getRoles(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All Roles invoked...")
	data, err := svc.Store.GetAllRoles(s)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, role := range data {
			updateRole(s, svc, role)
		}
		roleList := config.RoleList{}
		roleList.Roles = data
		roleList.Structure = uicomponent.GetRoleGridViewColumns(s)
		httputils.ServeJSON(w, roleList)
	}
}

func (svc *Service) getRoleFields(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get Role Form Fields")
	data := &config.Role{}
	httputils.ServeJSON(w, uicomponent.GetRoleConfigFormatter(s, data))
}

func (svc *Service) getRoleByID(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get Role By Id... Id=%v", id)
	data, err := svc.Store.GetRoleByID(s, id)
	if err != nil {
		logutil.Errorf(s, "Get Role Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		updateRole(s, svc, data)
		httputils.ServeJSON(w, uicomponent.GetRoleConfigFormatter(s, data))
	}
}

func (svc *Service) UpsertRole(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Add / Update Role Invoked")
	var role config.Role
	utils.DecodeAndValidate(s, w, req, &role)
	if nil != s.Err {
		return
	}
	logutil.Debugf(s, "Role object - %v ", role)
	err := svc.Store.UpsertRole(s, &role)
	if err != nil {
		logutil.Errorf(s, "Upsert Role Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, cluster := range role.Clusters {
			eventObj := utils.GetEventObj(cluster.UUID, role.EntityName(), role.URL(), role.ID, 0,
				utils.HttpPost, role)
			if req.Method == utils.HttpPut {
				eventObj.Data.CppmID = svc.Store.GetRoleClusterCPPMID(role.ID, cluster.ID, role.TenantID)
				if eventObj.Data.CppmID != 0 {
					eventObj.Data.Method = utils.HttpPut
				}
			}
			eventObj.TenantID = role.TenantID
			eventByte, _ := json.Marshal(eventObj)
			logutil.Debugf(s, "Notify Data  - %s", string(eventByte))
			go svc.Store.Watcher.Notify("event", eventByte)
		}
		httputils.ServeJSON(w, role)
	}
}

func (svc *Service) DeleteRole(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Delete Role by ID... Id = %v", id)

	var cppmID int
	var clusterID int
	var uuid string
	role, _ := svc.Store.GetRoleByID(s, id)

	for _, cluster := range role.Clusters {
		uuid = cluster.UUID
		clusterID = cluster.ID
		cppmID = svc.Store.GetRoleClusterCPPMID(role.ID, clusterID, role.TenantID)
		logutil.Debugf(s, "role ID - %d", role.ID)
		logutil.Debugf(s, "CPPM ID - %d", cppmID)
		if cppmID != 0 {
			eventObj := utils.GetEventObj(uuid, role.EntityName(), role.URL(), role.ID, cppmID,
				utils.HttpDelete, nil)

			logutil.Debugf(s, "Notify Data  - %v", eventObj)
			go svc.Store.Watcher.Notify("event", eventObj)
		}
	}

	err := svc.Store.DeleteRole(s, id)
	if err != nil {
		logutil.Errorf(s, "Delete Role Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func updateRole(s *model.SessionContext, svc *Service, data *config.Role) {
	data.AddedAtEpoc = getEpoc(data.AddedAt)
	data.UpdatedAtEpoc = getEpoc(data.UpdatedAt)
	if len(data.Clusters) > 0 {
		for _, cluster := range data.Clusters {
			updateCluster(s, svc, cluster)
		}
	}
}
