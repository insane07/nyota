package requestinterceptor

import (
	"nyota/backend/i18n"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	rstore "gopkg.in/boj/redistore.v1"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	// TODO: Fix the key string.
	key   = []byte("SfTuB!#&D55663*^329")
	store sessions.Store

	// Cokie Name
	loginCokieName = "auth-token-cookie"

	// User info key in session
	userAuthenticated = "authenticated"
	userNameKey       = "loggedin-user-name"
	userTenantIDKey   = "loggedin-user-tenant-id"
	userPermissionKey = "loggedin-user-permission"

	// 15 mins * 60 sec
	maxAge = 900
)

func init() {
	redisOpts := utils.RedisOptions("")
	rstore, storeErr := rstore.NewRediStore(10, redisOpts.Network, redisOpts.Addr, "", key)
	if storeErr != nil {
		logutil.Printf(nil, "Cookie store in use for sessions")
		// Case when session is stored in cookie in local disc.
		cookieStore := sessions.NewCookieStore(key)
		cookieStore.MaxAge(maxAge)
		store = cookieStore
	} else {
		logutil.Printf(nil, "Redis store in use for sessions")
		// Case when session is stored in cookie in local disc.
		rstore.SetMaxAge(maxAge)
		store = rstore
	}
}

/*StartSession should be called as part of login to create new session and save user details.*/
func StartSession(s *model.SessionContext, r *http.Request, w http.ResponseWriter,
	userName string, tenantID string, permission map[string]string) {

	session, _ := store.Get(r, loginCokieName)
	tempPerm, _ := json.Marshal(permission)

	// Set user as authenticated
	session.Values[userAuthenticated] = true
	session.Values[userNameKey] = userName
	session.Values[userTenantIDKey] = tenantID
	session.Values[userPermissionKey] = tempPerm
	session.Save(r, w)
	setUserContextDataForAPI(s, tenantID, userName, permission, r.Header.Get(utils.HTTPAcceptLanguageKey))
}

/*EndSession should be called as part of logout to clear session.*/
func EndSession(r *http.Request, w http.ResponseWriter) {

	session, _ := store.Get(r, loginCokieName)

	// Revoke users authentication
	session.Values[userAuthenticated] = false
	session.Values[userNameKey] = ""
	session.Values[userTenantIDKey] = ""
	session.Values[userPermissionKey] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)
}

func isUserLoggedIn(session *sessions.Session) bool {

	// Check if user is authenticated
	if auth, ok := session.Values[userAuthenticated].(bool); !ok || !auth {
		return false
	}
	return true
}

/*ValidateSession will verify User has the right session...*/
func ValidateSession() Interceptor {

	return func(f PrizmHandler) PrizmHandler {

		return func(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

			// Handle Session Here...
			session, _ := store.Get(r, loginCokieName)

			if !isUserLoggedIn(session) {
				logutil.Errorf(nil, "Session check failed. URL - %s  Method - %s ", r.URL, r.Method)
				//http.Error(w, "Forbidden: Access is denied", http.StatusForbidden)
				s.Err = &model.AppError{Type: utils.SessionError, Message: "Unauthorized", Code: http.StatusUnauthorized}
				return
			}

			// Add user and tenant info in request here...
			tenantID, _ := session.Values[userTenantIDKey].(string)
			userName, _ := session.Values[userNameKey].(string)
			permission, _ := session.Values[userPermissionKey].([]byte)
			lang := r.Header.Get(utils.HTTPAcceptLanguageKey)

			var permissionMap map[string]string
			err := json.Unmarshal(permission, &permissionMap)
			if err != nil {
				logutil.Errorf(s, "Failed to fetch permission from session... Assigning Analyst permission.")
				permissionMap = utils.AnalystUserRolePermission
			}
			setUserContextDataForAPI(s, tenantID, userName, permissionMap, lang)

			// Update max age...
			session.Options.MaxAge = maxAge
			session.Save(r, w)

			// Call the next handler in chain
			logutil.Debugf(s, "Session check passed for URL - %s", r.URL)
			f(s, w, r)
		}
	}
}

func setUserContextDataForAPI(s *model.SessionContext, tenantID string, userName string,
	permission map[string]string, lang string) {
	// Add user and tenant info in request here...
	u := s.User
	u.TenantId = tenantID
	u.UserName = userName
	u.Permission = permission
	s.Lang = lang
	s.TFunc = i18n.Translate(s)
	s.Err = nil
}
