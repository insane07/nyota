package api

import (
	"goprizm/httputils"
	"net/http"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/utils"

	"github.com/gorilla/mux"
)

func (svc *Service) getEvents(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Get All Events invoked...")
	data, err := svc.Store.GetAllEvents(s)
	if err != nil {
		logutil.Errorf(s, "Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		for _, event := range data {
			updateEvent(s, svc, event)
		}
		eventList := config.EventList{}
		eventList.Events = data
		httputils.ServeJSON(w, eventList)
	}
}

func (svc *Service) getEventByID(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get Event By Id... Id=%v", id)
	data, err := svc.Store.GetEventByID(s, id)
	if err != nil {
		logutil.Errorf(s, "Get Event Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		updateEvent(s, svc, data)
		httputils.ServeJSON(w, data)
	}
}

func (svc *Service) UpsertEvent(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	logutil.Debugf(s, "Service layer - Add / Update Event Invoked")
	var event config.Event
	utils.DecodeAndValidate(s, w, req, &event)
	if nil != s.Err {
		return
	}
	logutil.Debugf(s, "Event object - %v ", event)
	err := svc.Store.UpsertEvent(s, &event)
	if err != nil {
		logutil.Errorf(s, "Upsert Event Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		httputils.ServeJSON(w, event)
	}
}

func (svc *Service) DeleteEvent(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Delete Event by ID... Id = %v", id)
	err := svc.Store.DeleteEvent(s, id)
	if err != nil {
		logutil.Errorf(s, "Delete Event Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func updateEvent(s *model.SessionContext, svc *Service, data *config.Event) {
	data.AddedAtEpoc = getEpoc(data.AddedAt)
	data.UpdatedAtEpoc = getEpoc(data.UpdatedAt)
}
