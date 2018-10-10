package utils

const (
	// HTTP METHODS
	HttpPost   = "POST"
	HttpPut    = "PUT"
	HttpGet    = "GET"
	HttpDelete = "DELETE"

	// Supported Roles
	AdminUserRole   = "ADMIN"
	AnalystUserRole = "ANALYST"

	// UI Menu Permission Key -> API grouping
	AssetMenuPermissionKey                = "DEVICES"
	CustomClassificationMenuPermissionKey = "USER-CLASSIFIED-DEVICES"
	UnclassifiedMenuPermissionKey         = "UNCLASSIFIED-DEVICES"
	PolicyManagerMenuPermissionKey        = "DISCOVERY-SETTINGS"
	GenericMenuPermissionKey              = "COMMON-ASSET"

	// API Method Permissions Supported
	ModifyPermission = "MODIFY"
	ReadPermission   = "READ"
	BlockPermission  = "BLOCK"
)

var (
	// Supported Permissions
	AdminUserRolePermission   = updateUserAdminPermission()
	AnalystUserRolePermission = updateUserAnalystPermission()
)

func updateUserAdminPermission() map[string]string {
	adminPermission := make(map[string]string)
	adminPermission[AssetMenuPermissionKey] = ModifyPermission
	adminPermission[CustomClassificationMenuPermissionKey] = ModifyPermission
	adminPermission[UnclassifiedMenuPermissionKey] = ModifyPermission
	adminPermission[PolicyManagerMenuPermissionKey] = ModifyPermission
	return adminPermission
}

func updateUserAnalystPermission() map[string]string {
	analystPermission := make(map[string]string)
	analystPermission[AssetMenuPermissionKey] = ModifyPermission
	analystPermission[CustomClassificationMenuPermissionKey] = BlockPermission
	analystPermission[UnclassifiedMenuPermissionKey] = BlockPermission
	analystPermission[PolicyManagerMenuPermissionKey] = BlockPermission
	return analystPermission
}
