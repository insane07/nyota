package api

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/store"
	"nyota/backend/uicomponent"
	"nyota/backend/utils"
	"goprizm/httputils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (svc *Service) getClusters(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All Clusters invoked...")
	data, err := getAllClusters(svc.Store, s)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, cluster := range data {
			updateCluster(s, svc, cluster)
		}
		httputils.ServeJSON(w, data)
	}
}

func getAllClusters(store *store.Store, s *model.SessionContext) ([]*config.Cluster, error) {
	return store.GetClusters(s)
}

func (svc *Service) getClusterByID(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get Cluster By Id... Id=%v", id)
	data, err := svc.Store.GetClusterById(s, id)
	if err != nil {
		logutil.Errorf(s, "Get Cluster Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		updateCluster(s, svc, data)
		httputils.ServeJSON(w, uicomponent.GetClusterConfigFormatter(s, data))
	}
}

func updateCluster(s *model.SessionContext, svc *Service, data *config.Cluster) {
	data.AddedAtEpoc = getEpoc(data.AddedAt)
	data.UpdatedAtEpoc = getEpoc(data.UpdatedAt)

	cppmNodes, _ := svc.Store.GetCPPMNodesForCluster(s, strconv.Itoa(data.ID))
	if nil == cppmNodes {
		cppmNodes = make([]*config.CppmNode, 0)
	}
	data.CPPMNodes = cppmNodes
}

func (svc *Service) getClusterFields(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get Cluster Form Fields")
	data := &config.Cluster{}
	httputils.ServeJSON(w, uicomponent.GetClusterConfigFormatter(s, data))
}

func (svc *Service) UpsertCluster(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Add / Update Cluster Invoked")
	var cluster config.Cluster
	utils.DecodeAndValidate(s, w, req, &cluster)
	if nil != s.Err {
		return
	}
	logutil.Debugf(s, "Cluster object - %v ", cluster)
	_, err := svc.Store.UpsertCluster(s, &cluster)
	if err != nil {
		logutil.Errorf(s, "Upsert Cluster Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		//go svc.Store.Watcher.Notify("event", model.Event{cluster.ID, "ccc_cluster"})
		w.WriteHeader(http.StatusCreated)
	}
}

func (svc *Service) DeleteCluster(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Delete Cluster by id... Id=", id)
	err := svc.Store.DeleteClusterById(s, id)
	if err != nil {
		logutil.Errorf(s, "Delete Cluster Error - ", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func getEpoc(time time.Time) int64 {
	tm := time.Unix()
	if tm > 0 {
		return tm
	}
	return 0
}
