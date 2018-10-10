package store

import (
	"fmt"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"strings"
	"time"

	gorp "gopkg.in/gorp.v2"
)

//GetAllRoles - get all roles with Role details
func (store *Store) GetAllRoles(s *model.SessionContext) ([]*config.Role, error) {
	logutil.Debugf(s, "Store Layer - Get All Roles")
	var roles []*config.Role
	err := store.DB().Select(&roles, "Select * from CCC_Role where tenant_id=$1", s.User.TenantId)
	if err != nil {
		return nil, err
	}

	// for _, role := range roles {
	// 	var clusters []*config.Cluster
	// 	err = store.DB().Select(&clusters, `SELECT cluster.* FROM ccc_cluster cluster
	// 		JOIN ccc_role_cluster role_cluster on cluster.id = role_cluster.cluster_id
	// 		WHERE role_cluster.role_id = $1 and role_cluster.tenant_id = $2`, role.ID, role.TenantID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	role.Clusters = clusters
	// }
	return roles, nil
}

//GetRoleByID - get role based on id
func (store *Store) GetRoleByID(s *model.SessionContext, id string) (*config.Role, error) {
	logutil.Debugf(s, "Store Layer - Get Role By Id")
	var role *config.Role
	err := store.DB().SelectOne(&role, "Select * from CCC_Role where id=$1", id)
	if err != nil {
		return nil, err
	}

	var clusters []*config.Cluster
	err = store.DB().Select(&clusters, `SELECT cluster.* FROM ccc_cluster cluster 
		JOIN ccc_role_cluster role_cluster on cluster.id = role_cluster.cluster_id 
		WHERE role_cluster.role_id = $1 and role_cluster.tenant_id = $2`, role.ID, role.TenantID)
	if err != nil {
		return nil, err
	}
	role.Clusters = clusters

	return role, nil
}

//UpsertRole - insert or update role
func (store *Store) UpsertRole(s *model.SessionContext, role *config.Role) error {
	logutil.Debugf(s, "Store Layer - Upsert All Roles")
	err := execTx(s, store.DB(), func(tx *gorp.Transaction) (err error) {

		// upsert segment
		if role.ID == 0 {
			role.AddedAt = time.Now()
			role.UpdatedAt = role.AddedAt
			err = tx.Insert(role)
		} else {
			role.UpdatedAt = time.Now()
			_, err = tx.Update(role)
		}

		if err != nil {
			logutil.Errorf(s, "upsert role:(%s) failed: %v", role.Name, err)
			return err
		}

		//get existing clusters for role
		var existingRoleClusters []*config.RoleCluster
		err = store.DB().Select(&existingRoleClusters, "select * from ccc_role_cluster where role_id=$1", role.ID)
		if nil != err {
			logutil.Errorf(s, "Error in fetching clusters from role cluster association table,err:%v", err)
		}

		var createRoleClusterIDs []int
		var deleteRoleClusterIDs []int

		/*
			for _, cluster := range role.Clusters {
					if cluster.ID == 0 {
						cluster.AddedAt = time.Now()
						cluster.UpdatedAt = cluster.AddedAt
						err = tx.Insert(cluster)
					} else {
						cluster.UpdatedAt = time.Now()
						_, err = tx.Update(cluster)
					}

					if err != nil {
						logutil.Errorf(s, "upsert cluster:(%s) failed, error : %v", cluster.Name, err)
						return err
					}
			}
		*/
		for _, existingRoleClusters := range existingRoleClusters {
			isPresent := false
			for _, cluster := range role.Clusters {
				if existingRoleClusters.ClusterID == cluster.ID {
					isPresent = true
					break
				}
			}

			if !isPresent {
				deleteRoleClusterIDs = append(deleteRoleClusterIDs, existingRoleClusters.ClusterID)
			}
		}

		for _, cluster := range role.Clusters {
			isPresent := false
			for _, existingRoleClusters := range existingRoleClusters {
				if existingRoleClusters.ClusterID == cluster.ID {
					isPresent = true
					break
				}
			}
			if !isPresent {
				createRoleClusterIDs = append(createRoleClusterIDs, cluster.ID)
			}
		}

		//create new association
		var createRoleClusters []interface{}
		for _, clusterID := range createRoleClusterIDs {
			roleCluster := &config.RoleCluster{}
			roleCluster.ClusterID = clusterID
			roleCluster.RoleID = role.ID
			roleCluster.TenantID = role.TenantID
			createRoleClusters = append(createRoleClusters, roleCluster)
		}

		if len(createRoleClusters) > 0 {
			err = tx.Insert(createRoleClusters...)
			if err != nil {
				logutil.Errorf(s, "insert role cluster mappings failed: %v", err)
				return err
			}
		}

		// delete mappings
		if len(deleteRoleClusterIDs) > 0 {
			clusterIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(deleteRoleClusterIDs)), ","), "[]")
			_, err = tx.Exec("DELETE FROM ccc_role_cluster WHERE tenant_id = $1 and role_id = $2 and cluster_id in($3)", role.TenantID, role.ID, clusterIDs)
			if err != nil {
				logutil.Errorf(s, "delete role cluster mappings failed: %v", err)
				return err
			}
		}

		logutil.Debugf(s, "Upsert Role Successful")
		return nil
	})

	return err
}

func (store *Store) DeleteRole(s *model.SessionContext, id string) error {
	logutil.Debugf(s, "Store Layer - Delete Role By Id")
	_, errDelRoleCluster := store.DB().Exec("DELETE FROM CCC_ROLE_CLUSTER WHERE ROLE_ID = $1 and TENANT_ID = $2", id, s.User.TenantId)
	if errDelRoleCluster != nil {
		logutil.Debugf(s, "Deletion failed in mapping table of Role Cluster.")
		return errDelRoleCluster
	}
	_, errDelRole := store.DB().Exec("DELETE FROM CCC_ROLE WHERE ID = $1", id)
	if errDelRole != nil {
		logutil.Errorf(s, "Role deletion failed.")
		return errDelRole
	}
	return nil
}

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

func (store *Store) GetRoleClusterCPPMID(roleID int, clusterID int, tenantID string) int {
	rows, err := store.DB().Query(`SELECT cppm_id from ccc_role_cluster WHERE role_id = $1 and cluster_id=$2 and tenant_id = $3`, roleID, clusterID, tenantID)
	if err != nil {
		return 0
	}

	var cppmID int
	for rows.Next() {
		rows.Scan(&cppmID)
	}
	return cppmID
}

func (store *Store) UpdateRoleWithCPPMID(s *model.SessionContext, roleID int, uuid string, cppmID int) {
	res, err := store.DB().Exec("UPDATE CCC_ROLE_CLUSTER SET CPPM_ID=$1 WHERE ROLE_ID = $2 AND CLUSTER_ID = (SELECT ID FROM CCC_CLUSTER WHERE UUID=$3)", cppmID, roleID, uuid)
	if nil != err {
		logutil.Errorf(s, "CPPM ID updation failed in role cluster association table")
	}
	logutil.Debugf(s, "Result - %v", res)
}
