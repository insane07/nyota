package api

import (
	"encoding/json"
	"goprizm/httputils"
	"net/http"
	"nyota/backend/api/requestinterceptor"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/store"
	"nyota/backend/utils"
)

var (
	// Supported Permissions
	adminPermission   = make(map[string]string)
	analystPermission = make(map[string]string)

	users = model.UserTenantDetailsArray{
		model.UserTenantDetails{UserName: "admin@nyota.com", Password: "pass123", TenantID: "1",
			UserTenantAttributes: model.UserTenantAttributes{Role: utils.AdminUserRole, Permissions: utils.AdminUserRolePermission}},
		model.UserTenantDetails{UserName: "admin1@nyota.com", Password: "pass123", TenantID: "2",
			UserTenantAttributes: model.UserTenantAttributes{Role: utils.AdminUserRole, Permissions: utils.AdminUserRolePermission}},
	}
)

func (svc *Service) addRecords(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

	logutil.Printf(s, "Add Records Request Start...")

	for _, localUser := range users {
		svc.Store.UpsertUser(nil, &localUser)
	}
}

func (svc *Service) login(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

	logutil.Printf(s, "Authentication Request Start...")

	// Auth HERE....
	user := model.UserLogin{}
	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		logutil.Errorf(s, error.Error())
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	sysUser, loginSuccess := checkDbUser(s, user, svc.Store)
	if loginSuccess == false {
		sysUser, loginSuccess = checkLocalUser(user)
	}

	if loginSuccess == false {
		logutil.Printf(s, "Login Failed...")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	requestinterceptor.StartSession(s, r, w, sysUser.UserName, sysUser.TenantID, sysUser.UserTenantAttributes.Permissions)
	logutil.Printf(s, "Authentication Request Complete - "+
		"[Role is - %s Permission - %v] ", sysUser.UserTenantAttributes.Role, sysUser.UserTenantAttributes.Permissions)
	userBasicDetails := model.UserTenantBasicDetails{UserName: sysUser.UserName,
		Role: sysUser.UserTenantAttributes.Role, Permission: sysUser.UserTenantAttributes.Permissions}
	httputils.ServeJSON(w, userBasicDetails)
}

func checkLocalUser(user model.UserLogin) (*model.UserTenantDetails, bool) {
	for _, localUser := range users {
		if user.UserName == localUser.UserName && user.Password == localUser.Password {
			return &localUser, true
		}
	}
	return nil, false
}

func checkDbUser(s *model.SessionContext, user model.UserLogin, store *store.Store) (*model.UserTenantDetails, bool) {
	dbUser, err := store.GetUserByName(s, user.UserName)
	if err != nil {
		logutil.Errorf(nil, "User fetch failed...: %v", err)
		return nil, false
	}

	if user.UserName == dbUser.UserName && user.Password == dbUser.Password {
		return dbUser, true
	}
	return nil, false
}

func logout(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {
	logutil.Printf(s, "Logout called...")
	requestinterceptor.EndSession(r, w)
}
