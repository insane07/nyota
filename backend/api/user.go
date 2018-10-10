package api

import (
	"goprizm/httputils"
	"net/http"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/utils"

	"github.com/gorilla/mux"
)

func (svc *Service) getUsers(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All Users...")
	data, err := svc.Store.GetAllUsers(s)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		userList := model.UserTenantDetailsArray{}
		for _, userTenantDetails := range data {
			userList = append(userList, *userTenantDetails)
		}
		httputils.ServeJSON(w, userList)
	}
}

func (svc *Service) getUserByName(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	userName := mux.Vars(req)["userName"]
	logutil.Debugf(s, "Service layer - Get User By UserName... UserName=%v", userName)
	data, err := svc.Store.GetUserByName(s, userName)
	if err != nil {
		logutil.Errorf(s, "Get User Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		httputils.ServeJSON(w, data)
	}
}

// func (svc *Service) UpsertUser(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
// 	logutil.Debugf(s, "Service layer - Add / Update Role Invoked")
// 	var user model.UserTenantDetails
// 	utils.DecodeAndValidate(s, w, req, &user)
// 	if nil != s.Err {
// 		return
// 	}
// 	logutil.Debugf(s, "Role object - %v ", role)
// 	err := svc.Store.UpsertRole(s, &role)
// 	if err != nil {
// 		logutil.Errorf(s, "Upsert Role Error - %v", err)
// 		utils.SetSomethingWrong(s)
// 	} else {
// 		for _, cluster := range role.Clusters {
// 			eventObj := utils.GetEventObj(cluster.UUID, role.EntityName(), role.URL(), role.ID, 0,
// 				utils.HttpPost, role)
// 			if req.Method == utils.HttpPut {
// 				eventObj.Data.CppmID = svc.Store.GetRoleClusterCPPMID(role.ID, cluster.ID, role.TenantID)
// 				if eventObj.Data.CppmID != 0 {
// 					eventObj.Data.Method = utils.HttpPut
// 				}
// 			}
// 			eventObj.TenantID = role.TenantID
// 			eventByte, _ := json.Marshal(eventObj)
// 			logutil.Debugf(s, "Notify Data  - %s", string(eventByte))
// 			go svc.Store.Watcher.Notify("event", eventByte)
// 		}
// 		httputils.ServeJSON(w, role)
// 	}
// }
