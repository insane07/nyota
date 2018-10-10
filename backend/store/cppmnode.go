package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"time"
)

func (store *Store) GetCPPMNodes(s *model.SessionContext) ([]*config.CppmNode, error) {
	logutil.Debugf(s, "Store Layer - Get All CPPM Nodes")
	var cppmNodes []*config.CppmNode
	err := store.DB().Select(&cppmNodes, "SELECT * FROM CCC_CPPM_NODE WHERE TENANT_ID = $1", s.User.TenantId)
	if err != nil {
		return nil, err
	}
	return cppmNodes, nil
}

func (store *Store) GetCPPMNodesForCluster(s *model.SessionContext, clusterId string) ([]*config.CppmNode, error) {
	logutil.Debugf(s, "Store Layer - Get All CPPM Nodes By Cluster ID")
	var cppmNodes []*config.CppmNode
	err := store.DB().Select(&cppmNodes, "SELECT * FROM CCC_CPPM_NODE WHERE CLUSTER_ID = $1", clusterId)
	if err != nil {
		return nil, err
	}
	return cppmNodes, nil
}

func (store *Store) GetCPPMNodeById(s *model.SessionContext, id string) (*config.CppmNode, error) {
	logutil.Debugf(s, "Store Layer - Get CPPM Node By Id")
	var cppmNode *config.CppmNode
	err := store.DB().SelectOne(&cppmNode, "SELECT * FROM CCC_CPPM_NODE WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}
	return cppmNode, nil
}

func (store *Store) UpsertCPPMNode(s *model.SessionContext, data *config.CppmNode) (*config.CppmNode, error) {
	logutil.Debugf(s, "Store Layer - Upsert CPPM Node")

	if data.ID == 0 {
		data.AddedAt = time.Now()
		data.UpdatedAt = data.AddedAt
		err := store.DB().Insert(data)
		if err != nil {
			return nil, err
		}
	} else {
		data.UpdatedAt = time.Now()
		_, err := store.DB().Update(data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (store *Store) UpsertCPPMNodeEvent(s *model.SessionContext, data *config.CppmNode) error {
	logutil.Debugf(s, "Store Layer - Upsert CPPM Node")
	var cppmNode *config.CppmNode
	store.DB().SelectOne(&cppmNode, "SELECT * FROM CCC_CPPM_NODE WHERE SERVER_UUID = $1", data.ServerUUID)

	if nil != cppmNode {
		data.ID = cppmNode.ID
		data.UpdatedAt = time.Now()
		_, err := store.DB().Update(data)
		if err != nil {
			logutil.Errorf(s, "error in cluster update:%v", err)
			return err
		}
	} else {
		data.AddedAt = time.Now()
		data.UpdatedAt = data.AddedAt
		err := store.DB().Insert(data)
		if err != nil {
			logutil.Errorf(s, "error in CPPM Node insert:%v", err)
			return err
		}
	}
	return nil
}

func (store *Store) DeleteCPPMNode(s *model.SessionContext, id string) error {
	logutil.Debugf(s, "Store Layer - Delete CPPM Node By Id")
	_, err := store.DB().Exec("DELETE FROM CCC_CPPM_NODE WHERE ID = $1", id)
	if err != nil {
		return err
	}
	return nil
}

//GetClusterByUUID - fetches cluster based on uuid
func (store *Store) GetClusterByUUID(s *model.SessionContext, uuid string) *config.Cluster {
	var cluster *config.Cluster
	store.DB().SelectOne(&cluster, "SELECT * FROM CCC_CLUSTER WHERE UUID = $1", uuid)
	logutil.Debugf(s, "cluster object :%v", cluster)
	return cluster
}

//CreateAndFetchCluster = create a new one and send it back
func (store *Store) CreateAndFetchCluster(s *model.SessionContext, event model.Event) *config.Cluster {
	cluster := &config.Cluster{}
	cluster.UUID = event.UUID
	cluster.Name = event.UUID
	cluster.CppmVersion = event.CPPMVersion
	cluster.TenantID = event.TenantID
	cluster.AddedAt = time.Now()
	cluster.UpdatedAt = cluster.AddedAt
	err := store.DB().Insert(cluster)
	if err != nil {
		logutil.Errorf(s, "error in cluster insert:%v", err)
	}
	return cluster
}
