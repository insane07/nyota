package api

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/uicomponent"
	"nyota/backend/utils"
	"encoding/json"
	"goprizm/httputils"
	"net/http"

	"github.com/gorilla/mux"
)

func (svc *Service) getCPPMNodes(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All CPPM Nodes invoked...")
	data, err := svc.Store.GetCPPMNodes(s)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, cppmNode := range data {
			updateCppmNode(cppmNode)
		}
		httputils.ServeJSON(w, data)
	}
}

func (svc *Service) getCPPMNodesForCluster(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	clusterID := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get All CPPM Nodes by Cluster ID... Id = %v", clusterID)
	data, err := svc.Store.GetCPPMNodesForCluster(s, clusterID)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, cppmNode := range data {
			updateCppmNode(cppmNode)
		}
		httputils.ServeJSON(w, data)
	}
}

func (svc *Service) getCPPMNodeById(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get CPPM Node By ID... Id = %v", id)
	data, err := svc.Store.GetCPPMNodeById(s, id)
	if err != nil {
		logutil.Errorf(s, "Get CPPM Node Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		clusters, cerr := getAllClusters(svc.Store, s)
		if cerr != nil {
			logutil.Errorf(s, "Get All Clusters Error - %v", cerr)
			utils.SetSomethingWrong(s)
		} else {
			updateCppmNode(data)
			httputils.ServeJSON(w, uicomponent.GetCPPMNodeConfigFormatter(s, data, clusters))
		}
	}
}

func updateCppmNode(data *config.CppmNode) {
	data.AddedAtEpoc = getEpoc(data.AddedAt)
	data.UpdatedAtEpoc = getEpoc(data.UpdatedAt)
}

func (svc *Service) getCPPMNodeFields(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get CPPM Node Form Fields")
	data := &config.CppmNode{}
	clusters, cerr := getAllClusters(svc.Store, s)
	if cerr != nil {
		logutil.Errorf(s, "Get All Clusters Error - %v", cerr)
		utils.SetSomethingWrong(s)
	} else {
		httputils.ServeJSON(w, uicomponent.GetCPPMNodeConfigFormatter(s, data, clusters))
	}
}

func (svc *Service) UpsertCPPMNode(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Add / Update CPPM Node Invoked")
	var cppmNode config.CppmNode
	utils.DecodeAndValidate(s, w, req, &cppmNode)
	if nil != s.Err {
		return
	}
	logutil.Debugf(s, "CPPM Node object - %v ", cppmNode)
	_, err := svc.Store.UpsertCPPMNode(s, &cppmNode)
	if err != nil {
		logutil.Errorf(s, "Upsert CPPM Node Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (svc *Service) DeleteCPPMNode(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Delete CPPM Node by ID... Id = %v", id)
	err := svc.Store.DeleteCPPMNode(s, id)
	if err != nil {
		logutil.Errorf(s, "Delete CPPM Node Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (svc *Service) ExecuteEvent(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	event := model.Event{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&event)

	s = &model.SessionContext{}
	user := &model.UserContext{}
	user.TenantId = event.TenantID
	user.UserName = "abc"
	s.User = user

	logutil.Debugf(s, "Service layer - Add / Update CPPM Node Invoked")

	if err == nil {
		logutil.Debugf(s, "Entity Name:%s", event.Data.EntityName)
		switch event.Data.EntityName {
		case "Role":
			logutil.Debugf(s, "CCCID:%d", event.Data.CccID)
			logutil.Debugf(s, "CPPMID:%d", event.Data.CppmID)
			logutil.Debugf(s, "event:%v", event)
			if event.Data.CccID != 0 {
				svc.Store.UpdateRoleWithCPPMID(s, event.Data.CccID, event.UUID, event.Data.CppmID)
			}
			break
		case "cppmnode":
			uuid := event.UUID
			tenantID := event.TenantID
			var cppmNodes []config.CppmNode
			payload := event.Data.Payload
			payloadByte, _ := json.Marshal(payload)
			json.Unmarshal(payloadByte, &cppmNodes)
			logutil.Debugf(s, "CPPM Node object Array - %v ", cppmNodes)
			logutil.Debugf(s, "length of cppmnodes:%d", len(cppmNodes))

			//insert/update cluster table
			cluster := svc.Store.GetClusterByUUID(s, uuid)
			if nil == cluster {
				cluster = svc.Store.CreateAndFetchCluster(s, event)
			}

			for _, cppmNode := range cppmNodes {
				logutil.Debugf(s, "CPPM Node object - %v ", cppmNode)
				cppmNode.TenantID = tenantID
				cppmNode.ClusterID = cluster.ID
				logutil.Debugf(s, "new CPPM Node object - %v ", cppmNode)
				err := svc.Store.UpsertCPPMNodeEvent(s, &cppmNode)
				if err != nil {
					logutil.Errorf(s, "Upsert CPPM Node(%s) , error: %v", cppmNode.ServerIP, err)
				}
			}
			break
		}
	}
	event = model.Event{}
	httputils.ServeJSON(w, event)
}
