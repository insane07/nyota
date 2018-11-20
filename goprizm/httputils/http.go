package httputils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goprizm/log"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	// HTTP header Content-Type for JSON payload
	ContentTypeJSON = "application/json"
)

// IsLocalHostConn return true if HTTP request is send from inside localhost.
func IsLocalHostConn(req *http.Request) bool {
	remote := req.Header.Get("X-Forwarded-For")
	if remote == "" {
		remote = req.RemoteAddr
	}
	// Extract IP part if remote address is in the format ip:port
	parts := strings.SplitN(remote, ":", 2)
	remote = parts[0]

	return remote == "localhost" || remote == "127.0.0.1"
}

// Parse Authorization HTTP header and extract username and password.
func UserPasswordFromAutzHdr(autz string) (user string, password string, err error) {
	fields := strings.SplitN(autz, " ", 2)
	if len(fields) != 2 || fields[0] != "Basic" {
		return "", "", fmt.Errorf("expected `Basic Base64(user:password)`")
	}

	b64, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", "", fmt.Errorf("base64 decode user+password(%v)", err)
	}

	userPass := strings.SplitN(string(b64), ":", 2)
	if len(userPass) != 2 {
		return "", "", fmt.Errorf("expected user:password fmt")
	}

	return userPass[0], userPass[1], nil
}

// Find ip address of client given HTTP request. If X-FORWARDED-FOR header is set (proxy case), use it.
// Otherwise get ip:port from HTTP request socket.
func RemoteAddr(r *http.Request) string {
	remote := r.Header.Get("X-Forwarded-For")
	if remote == "" {
		remote = r.RemoteAddr
	}

	return remote
}

// Log error and send HTTP response.
func Error(w http.ResponseWriter, code int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Errorf(msg)
	http.Error(w, "", code)
}

// Send JSON HTTP response
func ServeJSON(w http.ResponseWriter, v interface{}) {
	ServeJSONWithStatus(w, v, http.StatusOK)
}

// ServeJSONWithStatus sends JSON HTTP response with specified HTTP status code
func ServeJSONWithStatus(w http.ResponseWriter, v interface{}, status int) {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(status)
	io.WriteString(w, string(content))
}

// Read HTTP request body as a JSON object
func ReadJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ReadError returns formatted error build from HTTP response status and body.
func ReadError(res *http.Response) error {
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return fmt.Errorf("HTTP status:%s (%s)", res.Status, string(body))
}
