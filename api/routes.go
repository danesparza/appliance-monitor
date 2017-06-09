package api

import (
	"encoding/json"
	"net/http"

	"time"
)

var (
	// BuildVersion contains the version information for the app
	BuildVersion = "Unknown"

	// CommitID is the git commitId for the app.  It's filled in as
	// part of the automated build
	CommitID string

	// ApplicationStartTime is the start time of the app
	ApplicationStartTime = time.Now()
)

// ErrorResponse represents an API response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//	Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := ErrorResponse{
		Status:  code,
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
