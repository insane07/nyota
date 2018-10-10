package utils

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

const (
	HTTPContentTypeKey   = "Content-Type"
	HTTPContentJSONValue = "application/json"

	HTTPAcceptKey   = "Accept"
	HTTPAcceptValue = "application/json"

	HTTPAcceptLanguageKey = "Accept-Language"

	notFoundError     = "Not Found Error"
	ValidatationError = "Validation Error"
	parsingError      = "Parsing Error"
	SessionError      = "Session Error"
	AccessError       = "Access Error"
	unkownError       = "Something Went Wrong"

	paramActiveFilter = "active_filter"
	ccParamFilter     = "Unclassified Device"
)

//DecodeAndValidate - Generic method which decodes the JSON body into model object and calls validate on the same
func DecodeAndValidate(s *model.SessionContext, w http.ResponseWriter, req *http.Request, m model.Context) {

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&m)
	if err != nil {
		logutil.Debugf(s, " Error %s", err.Error())
		SetParsingError(s, err)
		return
	}
	if req.Method != "POST" {
		m.SetData(mux.Vars(req)["id"], s.User.TenantId, s.User.UserName)
	} else {
		m.SetData("0", s.User.TenantId, s.User.UserName)
	}
	err = m.Validate()
	if nil != err {
		addValidationErrors(s, err)
		return
	}
	addAuditData(s, m)
}

func addAuditData(s *model.SessionContext, m model.Context) {
	s.AuditData = m.Audit()
}

func addValidationErrors(s *model.SessionContext, err error) {
	if nil != err {
		var data string
		switch err.(type) {
		case v.Errors:
			byteArr, _ := err.(v.Errors).MarshalJSON()
			m := make(map[string]interface{})
			json.Unmarshal(byteArr, &m)
			parseMap(s, m)
			d, _ := json.Marshal(m)
			data = string(d)
		default:
			// We should not enter here...
			data = err.Error()
		}
		s.Err = &model.AppError{Type: ValidatationError, Message: data, Code: http.StatusUnprocessableEntity}
	}
}
func parseMap(s *model.SessionContext, aMap map[string]interface{}) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(s, val.(map[string]interface{}))
		case []interface{}:
			parseArray(s, val.([]interface{}))
		case string:
			aMap[key] = s.TFunc(val.(string))
		default:
			logutil.Debugf(s, "Concrete Value :%s", concreteVal)
		}
	}
}

func parseArray(s *model.SessionContext, anArray []interface{}) {
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(s, val.(map[string]interface{}))
		case []interface{}:
			parseArray(s, val.([]interface{}))
		case string:
			val = s.TFunc(val.(string))
		default:
			logutil.Debugf(s, "Concrete Value :%s", concreteVal)
		}
	}
}

// SetPreconditionFailedError - Sets pre-condition Failed Error to session and handled generically.
func SetPreconditionFailedError(s *model.SessionContext, msg string) {
	s.Err = &model.AppError{Type: ValidatationError, Message: s.TFunc(msg), Code: http.StatusBadRequest}
}

// SetBadRequestError - Sets error to session and handled generically.
func SetBadRequestError(s *model.SessionContext) {
	s.Err = &model.AppError{Type: ValidatationError, Message: s.TFunc("name_unique_constraint_missing"), Code: http.StatusBadRequest}
}

// SetNotFoundError - Sets error to session and handled generically.
func SetNotFoundError(s *model.SessionContext) {
	s.Err = &model.AppError{Type: notFoundError, Message: "Not found.", Code: http.StatusNotFound}
}

// SetSomethingWrong - Sets error to session and handled generically if unintended error occurs.
func SetSomethingWrong(s *model.SessionContext) {
	s.Err = &model.AppError{Type: unkownError, Message: "Something went wrong...", Code: http.StatusInternalServerError}
}

// ReadHTTPResponse - read http response and return bytes
func ReadHTTPResponse(s *model.SessionContext, response *http.Response, err error) ([]byte, bool) {
	if err != nil {
		logutil.Errorf(s, "Generic HTTP Client error: %v", err)
		return nil, false
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logutil.Errorf(s, "Error reading the http response: %v", err)
		return nil, false
	}
	if response.StatusCode == http.StatusOK {
		return contents, true
	}
	logutil.Errorf(s, "Response status is not ok... Status: %d and  error: %s", response.StatusCode, string(contents))
	return contents, false
}

// ReadAndHandleHTTPResponse - read http response and return bytes
func ReadAndHandleHTTPResponse(s *model.SessionContext, response *http.Response, err error) ([]byte, bool) {
	if err != nil {
		logutil.Errorf(s, "Generic HTTP Client error: %v", err)
		return nil, false
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logutil.Errorf(s, "error reading the http response: %v", err)
		SetSomethingWrong(s)
		return nil, false
	}
	if response.StatusCode == http.StatusOK {
		return contents, true
	}
	logutil.Errorf(s, "Response status is not ok... Status: %d and  error: %s", response.StatusCode, string(contents))
	if response.StatusCode == http.StatusNotFound {
		SetNotFoundError(s)
	} else if response.StatusCode == http.StatusBadRequest {
		SetBadRequestError(s)
	} else {
		SetSomethingWrong(s)
	}
	return contents, false
}

// write byte array to response
func WriteHTTPResponse(contents []byte, w http.ResponseWriter) {
	w.Header().Set(HTTPContentTypeKey, HTTPContentJSONValue)
	w.Write(contents)
}

// SetParsingError - Error is set when unmarshalling fails.
func SetParsingError(s *model.SessionContext, err error) {
	s.Err = &model.AppError{Type: parsingError, Message: err.Error(), Code: http.StatusInternalServerError}
}

func IsValidIPV4(ip string) bool {
	IP := net.ParseIP(ip)
	return IP != nil && strings.Contains(ip, ".")
}

func ValidateIPs(value interface{}) error {
	ipArr := value.([]string)
	for _, ip := range ipArr {
		return ValidateIP(ip)
	}
	return nil
}

func ValidateIP(value interface{}) error {
	ip := value.(string)
	if !IsValidIPV4(ip) {
		return errors.New("key_invalid_ip")
	}
	return nil
}

// GetEventObj - Util method to prepare Redis Event Object
func GetEventObj(uuid string, name string, url string, cccID int, cppmID int,
	method string, data interface{}) model.Event {

	eventData := model.EventData{EntityName: name, CccID: cccID, CppmID: cppmID, Method: method,
		URI: url, Payload: data}
	if method == HttpPut || method == HttpDelete {
		eventData.URI = fmt.Sprintf("%s/%d", url, cppmID)
	}
	eventObj := model.Event{UUID: uuid, Data: eventData}
	return eventObj
}
