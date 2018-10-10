package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"time"

	gorp "gopkg.in/gorp.v2"
)

//GetAllEvents - get all events with Event details
func (store *Store) GetAllEvents(s *model.SessionContext) ([]*config.Event, error) {
	logutil.Debugf(s, "Store Layer - Get All Events")
	var events []*config.Event
	err := store.DB().Select(&events, "Select * from Events where tenant_id=$1", s.User.TenantId)
	if err != nil {
		return nil, err
	}

	// for _, event := range events {
	// 	var clusters []*config.Cluster
	// 	err = store.DB().Select(&clusters, `SELECT cluster.* FROM ccc_cluster cluster
	// 		JOIN ccc_event_cluster event_cluster on cluster.id = event_cluster.cluster_id
	// 		WHERE event_cluster.event_id = $1 and event_cluster.tenant_id = $2`, event.ID, event.TenantID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	event.Clusters = clusters
	// }
	return events, nil
}

//GetEventByID - get event based on id
func (store *Store) GetEventByID(s *model.SessionContext, id string) (*config.Event, error) {
	logutil.Debugf(s, "Store Layer - Get Event By Id")
	var event *config.Event
	err := store.DB().SelectOne(&event, "Select * from Events where id=$1", id)
	if err != nil {
		return nil, err
	}

	// var clusters []*config.Cluster
	// err = store.DB().Select(&clusters, `SELECT cluster.* FROM ccc_cluster cluster
	// 	JOIN ccc_event_cluster event_cluster on cluster.id = event_cluster.cluster_id
	// 	WHERE event_cluster.event_id = $1 and event_cluster.tenant_id = $2`, event.ID, event.TenantID)
	// if err != nil {
	// 	return nil, err
	// }
	// event.Clusters = clusters

	return event, nil
}

//UpsertEvent - insert or update event
func (store *Store) UpsertEvent(s *model.SessionContext, event *config.Event) error {
	logutil.Debugf(s, "Store Layer - Upsert All Events")
	err := execTx(s, store.DB(), func(tx *gorp.Transaction) (err error) {

		// upsert segment
		if event.ID == 0 {
			event.AddedAt = time.Now()
			event.UpdatedAt = event.AddedAt
			err = tx.Insert(event)
		} else {
			event.UpdatedAt = time.Now()
			_, err = tx.Update(event)
		}

		if err != nil {
			logutil.Errorf(s, "upsert event:(%s) failed: %v", event.Name, err)
			return err
		}
		logutil.Debugf(s, "Upsert Event Successful")
		return nil
	})

	return err
}

func (store *Store) DeleteEvent(s *model.SessionContext, id string) error {
	logutil.Debugf(s, "Store Layer - Delete Event By Id")
	_, errDelEventCluster := store.DB().Exec("DELETE FROM CCC_ROLE_CLUSTER WHERE ROLE_ID = $1 and TENANT_ID = $2", id, s.User.TenantId)
	if errDelEventCluster != nil {
		logutil.Debugf(s, "Deletion failed in mapping table of Event Cluster.")
		return errDelEventCluster
	}
	_, errDelEvent := store.DB().Exec("DELETE FROM CCC_ROLE WHERE ID = $1", id)
	if errDelEvent != nil {
		logutil.Errorf(s, "Event deletion failed.")
		return errDelEvent
	}
	return nil
}

/*
func (store *Store) GetEventCluster(eventID int, clusterID int, tenantID string) *config.EventCluster {
	var eventCluster *config.EventCluster
	err := store.DB().Select(&eventCluster, "SELECT * from ccc_event_cluster WHERE event_id = $1 and cluster_id=$2 and tenant_id = $3", eventID, clusterID, tenantID)
	if err != nil {
		return nil
	}
	return eventCluster
}
*/

func (store *Store) GetEventClusterCPPMID(eventID int, clusterID int, tenantID string) int {
	rows, err := store.DB().Query(`SELECT cppm_id from ccc_event_cluster WHERE event_id = $1 and cluster_id=$2 and tenant_id = $3`, eventID, clusterID, tenantID)
	if err != nil {
		return 0
	}

	var cppmID int
	for rows.Next() {
		rows.Scan(&cppmID)
	}
	return cppmID
}

func (store *Store) UpdateEventWithCPPMID(s *model.SessionContext, eventID int, uuid string, cppmID int) {
	res, err := store.DB().Exec("UPDATE CCC_ROLE_CLUSTER SET CPPM_ID=$1 WHERE ROLE_ID = $2 AND CLUSTER_ID = (SELECT ID FROM CCC_CLUSTER WHERE UUID=$3)", cppmID, eventID, uuid)
	if nil != err {
		logutil.Errorf(s, "CPPM ID updation failed in event cluster association table")
	}
	logutil.Debugf(s, "Result - %v", res)
}
