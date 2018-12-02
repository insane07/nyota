package api

import (
	"bytes"
	"goprizm/httputils"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"nyota/backend/utils"
	"strconv"

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

func (svc *Service) getEventQrByID(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	logutil.Debugf(s, "Service layer - Get Event QR By Id... Id=%v", id)
	img, err := svc.Store.GetEventQrByID(s, id)
	if err != nil {
		logutil.Errorf(s, "Get Event QR Error - %v", err)
		utils.SetSomethingWrong(s)
	} else {
		writeImage(w, img)
	}
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
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
