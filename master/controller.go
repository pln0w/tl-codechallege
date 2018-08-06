package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Controller struct{}

// NewController - returns pointer to Controller struct
func NewController() *Controller {
	return &Controller{}
}

// SendJSON - function returns JSON response of any object
func (ctrl *Controller) SendJSON(w http.ResponseWriter, v interface{}, code int) {

	w.Header().Add("Content-Type", "application/json")

	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(b))
	}
}

// HandleError - function returns error JSON message
func (ctrl *Controller) HandleError(err error, w http.ResponseWriter, status ...int) {

	msg := map[string]string{
		"status":  "fail",
		"message": err.Error(),
	}

	returnStatus := http.StatusInternalServerError
	if len(status) > 0 {
		returnStatus = status[0]
	}

	ctrl.SendJSON(w, &msg, returnStatus)

}

// Index - function gets request, dispatch proper action to workers via gRPC
// and return results back received content
func (ctrl *Controller) Index(w http.ResponseWriter, r *http.Request) {

	// TODO: gRPC implementation

	res := map[string]interface{}{
		"status": http.StatusOK,
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}

// HealthCheck - function returns
func (ctrl *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		ctrl.HandleError(err, w, http.StatusInternalServerError)
	}

	res := map[string]interface{}{
		"status":   http.StatusOK,
		"hostname": hostname,
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}
