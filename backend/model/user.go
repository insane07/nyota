package model

import "github.com/nicksnyder/go-i18n/i18n"

type UserContext struct {
	TenantId   string
	UserName   string
	Permission map[string]string
}

type AppError struct {
	Type, Message string
	Code          int
}

type SessionContext struct {
	User      *UserContext
	Lang      string             // Client language preference
	TFunc     i18n.TranslateFunc // I18N Translation function based on client language preference
	Err       *AppError
	AuditData string
}

type UserLogin struct {
	UserName, Password string
}

type UserTenantBasicDetails struct {
	UserName   string            `json:"username"`
	Role       string            `json:"role"`
	Permission map[string]string `json:"permission"`
}

type UserTenantDetails struct {
	UserName             string `db:"username" json:"name"`
	Password             string `db:"password" json:"password"`
	TenantID             string `db:"tenant_id" json:"tenant_id"`
	Descrition           string `db:"description" json:"descrption"`
	UserTenantAttributes `json:"attributes"`
}

type UserTenantAttributes struct {
	Role        string            `json:"role"`
	Permissions map[string]string `json:"permissions"`
}

type UserTenantDetailsArray []UserTenantDetails

//Context interface. All the models which needs to be validated will implement this interface
type Context interface {
	Validate() error                                     // bool to identify add/edit operation
	Audit() string                                       // Audit data
	SetData(id string, tenantID string, userName string) // set Id to entity for put request
}
