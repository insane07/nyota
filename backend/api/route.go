package api

import (
	"nyota/backend/api/requestinterceptor"
	"nyota/backend/utils"
)

/*Route has details about he route*/
type Route struct {
	Path, Name         string
	Method, Permission string
	RealHandler        requestinterceptor.PrizmHandler
	Group              string
}

/*Routes defines all routes in the system*/
type Routes []Route

func getAllRoutes(srv *Service) (Routes, Routes) {

	nologinRoutes := Routes{
		Route{"/login", "Login", utils.HttpPost, utils.ReadPermission, srv.login, utils.GenericMenuPermissionKey},
		Route{"/init", "Init", utils.HttpPost, utils.ReadPermission, srv.addRecords, utils.GenericMenuPermissionKey},
		Route{"/event", "execute event", utils.HttpPost, utils.ModifyPermission, srv.ExecuteEvent, utils.GenericMenuPermissionKey},
	}
	/*GuardedRoutes are routes with Login*/
	guardedRoutes := Routes{
		Route{"/logout", "Logout", utils.HttpGet, utils.ReadPermission, logout, utils.GenericMenuPermissionKey},

		Route{"/events", "Get-Events", utils.HttpGet, utils.ReadPermission, srv.getEvents, utils.GenericMenuPermissionKey},
		Route{"/events/{id:[0-9]+}", "Get-Event-By-Id", utils.HttpGet, utils.ReadPermission, srv.getEventByID, utils.GenericMenuPermissionKey},
		Route{"/events", "Add-Event", utils.HttpPost, utils.ModifyPermission, srv.UpsertEvent, utils.GenericMenuPermissionKey},
		Route{"/events/{id:[0-9]+}", "Update-Event-By-Id", utils.HttpPut, utils.ModifyPermission, srv.UpsertEvent, utils.GenericMenuPermissionKey},
		Route{"/events/{id:[0-9]+}", "Delete-Event-By-Id", utils.HttpDelete, utils.ModifyPermission, srv.DeleteEvent, utils.GenericMenuPermissionKey},

		Route{"/tenants", "Get-Tenants", utils.HttpGet, utils.ReadPermission, srv.getTenants, utils.GenericMenuPermissionKey},
		Route{"/tenants/{id:[0-9]+}", "Get-Tenant-By-Id", utils.HttpGet, utils.ReadPermission, srv.getTenantById, utils.GenericMenuPermissionKey},
		Route{"/tenants", "Add-Tenant", utils.HttpPost, utils.ModifyPermission, srv.UpsertTenant, utils.GenericMenuPermissionKey},
		Route{"/tenants/{id:[0-9]+}", "Update-Tenant-By-Id", utils.HttpPut, utils.ModifyPermission, srv.UpsertTenant, utils.GenericMenuPermissionKey},
		Route{"/tenants/{id:[0-9]+}", "Delete-Tenant-By-Id", utils.HttpDelete, utils.ModifyPermission, srv.DeleteTenant, utils.GenericMenuPermissionKey},

		Route{"/clusters", "Get-clusters", utils.HttpGet, utils.ReadPermission, srv.getClusters, utils.GenericMenuPermissionKey},
		Route{"/clusters/{id:[0-9]+}", "Get-clusters-By-Id", utils.HttpGet, utils.ReadPermission, srv.getClusterByID, utils.GenericMenuPermissionKey},
		Route{"/clusters", "Add-clusters", utils.HttpPost, utils.ModifyPermission, srv.UpsertCluster, utils.GenericMenuPermissionKey},
		Route{"/clusters/{id:[0-9]+}", "Update-clusters-By-Id", utils.HttpPut, utils.ModifyPermission, srv.UpsertCluster, utils.GenericMenuPermissionKey},
		Route{"/clusters/{id:[0-9]+}", "Delete-clusters-By-Id", utils.HttpDelete, utils.ModifyPermission, srv.DeleteCluster, utils.GenericMenuPermissionKey},
		Route{"/clusters/formfields", "cluster fields", utils.HttpGet, utils.ReadPermission, srv.getClusterFields, utils.GenericMenuPermissionKey},

		Route{"/cppmnodes", "Get-CPPMNodes", utils.HttpGet, utils.ReadPermission, srv.getCPPMNodes, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes/cluster/{id:[0-9]+}", "Get-CPPMNodes-By-ClusterId", utils.HttpGet, utils.ReadPermission, srv.getCPPMNodesForCluster, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes/{id:[0-9]+}", "Get-CPPMNodes-By-Id", utils.HttpGet, utils.ReadPermission, srv.getCPPMNodeById, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes", "Add-CPPMNodes", utils.HttpPost, utils.ModifyPermission, srv.UpsertCPPMNode, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes/{id:[0-9]+}", "Update-CPPMNodes-By-Id", utils.HttpPut, utils.ModifyPermission, srv.UpsertCPPMNode, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes/{id:[0-9]+}", "Delete-CPPMNodes-By-Id", utils.HttpDelete, utils.ModifyPermission, srv.DeleteCPPMNode, utils.GenericMenuPermissionKey},
		Route{"/cppmnodes/formfields", "CPPMNodes fields", utils.HttpGet, utils.ReadPermission, srv.getCPPMNodeFields, utils.GenericMenuPermissionKey},

		Route{"/roles", "Get-Roles", utils.HttpGet, utils.ReadPermission, srv.getRoles, utils.GenericMenuPermissionKey},
		Route{"/roles/{id:[0-9]+}", "Get-Role-By-Id", utils.HttpGet, utils.ReadPermission, srv.getRoleByID, utils.GenericMenuPermissionKey},
		Route{"/roles/formfields", "role fields", utils.HttpGet, utils.ReadPermission, srv.getRoleFields, utils.GenericMenuPermissionKey},
		Route{"/roles", "Add-Role", utils.HttpPost, utils.ModifyPermission, srv.UpsertRole, utils.GenericMenuPermissionKey},
		Route{"/roles/{id:[0-9]+}", "Update-Role-By-Id", utils.HttpPut, utils.ModifyPermission, srv.UpsertRole, utils.GenericMenuPermissionKey},
		Route{"/roles/{id:[0-9]+}", "Delete-Role-By-Id", utils.HttpDelete, utils.ModifyPermission, srv.DeleteRole, utils.GenericMenuPermissionKey},
	}
	return nologinRoutes, guardedRoutes
}
