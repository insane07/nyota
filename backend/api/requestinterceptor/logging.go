package requestinterceptor

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/utils"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// TrackReqResp logs all requests with its path and the time it took to process
func TrackReqResp() Interceptor {

	// Create a new Middleware
	return func(f PrizmHandler) PrizmHandler {

		// Define the http.HandlerFunc
		return func(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			logutil.Debugf(s, "URL - %s  Method - %s  Started", r.URL, r.Method)

			// Generic recovery in case panic occurred and not handled
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 1<<16)
					stackSize := runtime.Stack(buf, true)

					logutil.Errorf(s, "URL - %s  Method - %s  Time Taken - %s \n*** Recovered from panic ***\n StackTrace = %s",
						r.URL, r.Method, time.Since(start), string(buf[0:stackSize]))

					utils.SetSomethingWrong(s)
					return
				}
				logutil.Debugf(s, "URL - %s  Method - %s Completed with Time Taken - %s", r.URL, r.Method, time.Since(start))

				addAuditLogs(s, r)
			}()

			// Call the next middleware/handler in chain
			f(s, w, r)
		}
	}
}

func addAuditLogs(s *model.SessionContext, r *http.Request) {
	// Add Audit logs only when there is no error message
	if s.Err == nil {
		if "" != s.AuditData {
			logutil.Debugf(s, "Audit Log - Entity:%s, Action:%s, Data:%s",
				getEntity(r.URL.Path), getAction(r.Method), s.AuditData)
		} else {
			logutil.Debugf(s, "Audit Log - Entity:%s, Action:%s",
				getEntity(r.URL.Path), getAction(r.Method))
		}

	}
}

func getEntity(a string) string {
	s := strings.SplitAfter(a, "/")
	return s[len(s)-1]
}

func getAction(m string) string {
	switch m {
	case "GET":
		return "Read"
	case "POST":
		return "Add"
	case "PUT":
		return "Update"
	case "DELETE":
		return "Delete"
	}
	return ""
}
