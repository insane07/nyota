package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"time"
)

func (store *Store) GetClusters(s *model.SessionContext) ([]*config.Cluster, error) {

	logutil.Debugf(s, "Store Layer - Get All Clusters")
	var clusters []*config.Cluster
	err := store.DB().Select(&clusters, "Select * from CCC_Cluster where tenant_id=$1", s.User.TenantId)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func (store *Store) GetClusterById(s *model.SessionContext, id string) (*config.Cluster, error) {

	logutil.Debugf(s, "Store Layer - Get Cluster By Id")
	var cluster *config.Cluster
	err := store.DB().SelectOne(&cluster, "Select * from CCC_Cluster where id=$1", id)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (store *Store) UpsertCluster(s *model.SessionContext, data *config.Cluster) (*config.Cluster, error) {

	logutil.Debugf(s, "Store Layer - Upsert Cluster")
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

func (store *Store) DeleteClusterById(s *model.SessionContext, id string) error {

	logutil.Debugf(s, "Store Layer - Delete Cluster By Id")
	_, err := store.DB().Exec("Delete from CCC_Cluster where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
