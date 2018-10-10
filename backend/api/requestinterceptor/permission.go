package requestinterceptor

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/utils"
	"net/http"
)

/*RBACCheck will verify user has right permission...*/
func RBACCheck(group, permission string) Interceptor {

	// Create a new Middleware
	return func(f PrizmHandler) PrizmHandler {

		// Define the http.HandlerFunc
		return func(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

			var isAllowed = false
			if group == utils.GenericMenuPermissionKey {
				// Common API no permission check...
				isAllowed = true
			} else {
				userPermission := s.User.Permission
				val, ok := userPermission[group]
				if ok {
					isAllowed = checkPermissions(s, permission, val)
				} else {
					logutil.Errorf(s, "RBAC check failed.."+
						"Unable to fetch permission for group - %s permission - %v", group, userPermission)
				}
			}
			// If  not allowed then throw HTTP 403
			if isAllowed == false {
				// Permission and group not matching...
				logutil.Errorf(s, "RBAC check failed for URL - %s", r.URL)
				s.Err = &model.AppError{Type: utils.AccessError, Message: "Forbidden: Access is denied", Code: http.StatusForbidden}
				return
			}
			// Call the next middleware/handler in chain
			logutil.Debugf(s, "RBAC check passed for URL - %s", r.URL)
			f(s, w, r)
		}
	}
}

/*
Blocked permission no access.
Matching permission allow access.
Higher permission allow access. "Modify vs Read" case
*/
func checkPermissions(s *model.SessionContext, methodPermission, userPermission string) bool {
	// Blocked permission no access.
	if userPermission == utils.BlockPermission {
		logutil.Errorf(s, "RBAC check failed... Permission is set to blocked...")
		return false
	}

	// Lower permission Deny access. "Read vs Modify" case
	if userPermission == utils.ReadPermission &&
		methodPermission == utils.ModifyPermission {
		logutil.Errorf(s, "RBAC check failed..."+
			"Method permission - %s & User Permission - %s ", methodPermission, userPermission)
		return false
	}

	// Matching permission allow access.
	// Higher permission allow access. "Modify vs Read" case
	return true
}
