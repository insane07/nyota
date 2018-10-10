package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"

	gorp "gopkg.in/gorp.v2"
)

//GetAllUsers - get all users with user details
func (store *Store) GetAllUsers(s *model.SessionContext) ([]*model.UserTenantDetails, error) {
	logutil.Debugf(s, "Store Layer - Get All Roles")
	var users []*model.UserTenantDetails
	err := store.DB().Select(&users, "Select * from USER_Tenant_Details where tenant_id=$1", s.User.TenantId)
	if err != nil {
		return nil, err
	}

	// for _, user := range users {
	// 	var attributes *model.UserTenantAttributes
	// 	err = store.DB().Select(&attributes, `SELECT cluster.* FROM ccc_cluster cluster
	// 		JOIN ccc_role_cluster role_cluster on cluster.id = role_cluster.cluster_id
	// 		WHERE role_cluster.role_id = $1 and role_cluster.tenant_id = $2`, role.ID, role.TenantID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	user.Attributes = attributes
	// }
	return users, nil
}

//GetUserByName - get user based on username
func (store *Store) GetUserByName(s *model.SessionContext, username string) (*model.UserTenantDetails, error) {
	logutil.Debugf(s, "Store Layer - Get User By UserName")
	var user *model.UserTenantDetails
	err := store.DB().SelectOne(&user, "Select * from user_tenant_details where username=$1", username)
	if err != nil {
		return nil, err
	}

	// var clusters []*config.Cluster
	// err = store.DB().Select(&clusters, `SELECT cluster.* FROM ccc_cluster cluster
	// 	JOIN ccc_role_cluster role_cluster on cluster.id = role_cluster.cluster_id
	// 	WHERE role_cluster.role_id = $1 and role_cluster.tenant_id = $2`, role.ID, role.TenantID)
	// if err != nil {
	// 	return nil, err
	// }
	// role.Clusters = clusters

	return user, nil
}

//UpsertUser - insert or update user
func (store *Store) UpsertUser(s *model.SessionContext, user *model.UserTenantDetails) error {
	logutil.Debugf(s, "Store Layer - Upsert All Users")
	err := execTx(s, store.DB(), func(tx *gorp.Transaction) (err error) {
		// upsert user
		err = tx.Insert(user)

		if err != nil {
			logutil.Errorf(s, "upsert user:(%s) failed: %v", user.UserName, err)
			return err
		}

		logutil.Debugf(s, "Upsert Role Successful")
		return nil
	})

	return err
}

// func (store *Store) DeleteRole(s *model.SessionContext, id string) error {
// 	logutil.Debugf(s, "Store Layer - Delete Role By Id")
// 	_, errDelRoleCluster := store.DB().Exec("DELETE FROM CCC_ROLE_CLUSTER WHERE ROLE_ID = $1 and TENANT_ID = $2", id, s.User.TenantId)
// 	if errDelRoleCluster != nil {
// 		logutil.Debugf(s, "Deletion failed in mapping table of Role Cluster.")
// 		return errDelRoleCluster
// 	}
// 	_, errDelRole := store.DB().Exec("DELETE FROM CCC_ROLE WHERE ID = $1", id)
// 	if errDelRole != nil {
// 		logutil.Errorf(s, "Role deletion failed.")
// 		return errDelRole
// 	}
// 	return nil
// }

/*
func (store *Store) GetRoleCluster(roleID int, clusterID int, tenantID string) *config.RoleCluster {
	var roleCluster *config.RoleCluster
	err := store.DB().Select(&roleCluster, "SELECT * from ccc_role_cluster WHERE role_id = $1 and cluster_id=$2 and tenant_id = $3", roleID, clusterID, tenantID)
	if err != nil {
		return nil
	}
	return roleCluster
}
*/

// func (store *Store) GetRoleClusterCPPMID(roleID int, clusterID int, tenantID string) int {
// 	rows, err := store.DB().Query(`SELECT cppm_id from ccc_role_cluster WHERE role_id = $1 and cluster_id=$2 and tenant_id = $3`, roleID, clusterID, tenantID)
// 	if err != nil {
// 		return 0
// 	}

// 	var cppmID int
// 	for rows.Next() {
// 		rows.Scan(&cppmID)
// 	}
// 	return cppmID
// }

// func (store *Store) UpdateRoleWithCPPMID(s *model.SessionContext, roleID int, uuid string, cppmID int) {
// 	res, err := store.DB().Exec("UPDATE CCC_ROLE_CLUSTER SET CPPM_ID=$1 WHERE ROLE_ID = $2 AND CLUSTER_ID = (SELECT ID FROM CCC_CLUSTER WHERE UUID=$3)", cppmID, roleID, uuid)
// 	if nil != err {
// 		logutil.Errorf(s, "CPPM ID updation failed in role cluster association table")
// 	}
// 	logutil.Debugf(s, "Result - %v", res)
// }
