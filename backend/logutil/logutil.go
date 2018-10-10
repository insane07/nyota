package logutil

import (
	"fmt"
	"goprizm/log"
	"nyota/backend/model"
)

const (
	adminbackend = "NYOTA"
)

//Debugf - debug log
func Debugf(s *model.SessionContext, format string, stringVal ...interface{}) {
	format = fmt.Sprintf("%s - %s", adminbackend, format)
	if nil == s {
		log.Debugf(format, stringVal...)
	} else {
		log.T(s.User.TenantId, "UserName", s.User.UserName).Debugf(format, stringVal...)
	}
}

//Errorf - error log
func Errorf(s *model.SessionContext, format string, stringVal ...interface{}) {
	format = fmt.Sprintf("%s - %s", adminbackend, format)
	if nil == s {
		log.Errorf(format, stringVal...)
	} else {
		log.T(s.User.TenantId, "UserName", s.User.UserName).Errorf(format, stringVal...)
	}
}

//Printf - print log
func Printf(s *model.SessionContext, format string, stringVal ...interface{}) {
	format = fmt.Sprintf("%s - %s", adminbackend, format)
	if nil == s {
		log.Printf(format, stringVal...)
	} else {
		log.T(s.User.TenantId, "UserName", s.User.UserName).Printf(format, stringVal...)
	}
}
