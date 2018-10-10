package main

import (
	"nyota/backend/api"
	"nyota/backend/logutil"
	"goprizm/log"
	"goprizm/sysutils"
	"net/http"
	"strconv"
)

var (
	adminPort = ":9100" // Backend Admin default port
)

func init() {
	port := sysutils.GetenvInt("ADMIN_BACKEND_PORT", 9100)

	if port > 1024 && port <= 65535 {
		adminPort = ":" + strconv.Itoa(port)
	}
	logutil.Printf(nil, "Admin Backend Service started on port%s", adminPort)
}

func main() {
	router := api.NewRoute()
	log.Fatalf("Error starting Admin Backend Service (%s)", http.ListenAndServe(adminPort, router))
}
