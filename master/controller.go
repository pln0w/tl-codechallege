package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// Controller structure
type Controller struct{}

// NewController - returns pointer to Controller struct
func NewController() *Controller {
	return &Controller{}
}

// SendJSON - function returns JSON response of any object
func (ctrl *Controller) SendJSON(w http.ResponseWriter, v interface{}, code int) {

	// Add proper content type header
	w.Header().Add("Content-Type", "application/json")

	// Marshal any object to JSON format
	content, err := json.Marshal(v)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {

		// Write response
		w.WriteHeader(code)
		io.WriteString(w, string(content))
	}
}

// HandleError - function returns error JSON message
func (ctrl *Controller) HandleError(err error, w http.ResponseWriter, status ...int) {

	// Prepare message
	msg := map[string]string{
		"status":  "fail",
		"message": err.Error(),
	}

	log.Error(err.Error())

	// Set proper return status or let 500 as default
	returnStatus := http.StatusInternalServerError
	if len(status) > 0 {
		returnStatus = status[0]
	}

	ctrl.SendJSON(w, &msg, returnStatus)
}

// Index - function gets request, dispatch proper action to workers via gRPC
// and return results back received content
func (ctrl *Controller) Index(w http.ResponseWriter, r *http.Request) {

	// Mock URL
	// url := os.Getenv("LB_HOST") + ":" + os.Getenv("LB_PORT") + "/"

	// Call any slave for results
	// response := CallSlave(url)

	response := "TODO: there is an issue while connecting other server"

	res := map[string]interface{}{
		"status":   http.StatusOK,
		"response": response,
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}

// HealthCheck - function returns status and hostname
func (ctrl *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {

	// Get hostname to be return back in response
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		ctrl.HandleError(err, w, http.StatusInternalServerError)
		return
	}

	res := map[string]interface{}{
		"status":   http.StatusOK,
		"hostname": hostname,
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}
