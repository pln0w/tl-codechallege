package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BaseController structure
type BaseController struct{}

// NewBaseController - returns pointer to BaseController struct
func NewBaseController() *BaseController {
	return &BaseController{}
}

// SendJSON - function returns JSON response of any object
func (ctrl *BaseController) SendJSON(w http.ResponseWriter, v interface{}, code int) {

	// Add proper content type header
	w.Header().Add("Content-Type", "application/json")

	// Marshal any object to JSON format
	content, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%v", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {

		// Write response
		w.WriteHeader(code)
		io.WriteString(w, string(content))
	}
}

// HandleError - function returns error JSON message
func (ctrl *BaseController) HandleError(err error, w http.ResponseWriter, status ...int) {

	// Prepare message
	msg := map[string]string{
		"status":  "fail",
		"message": err.Error(),
	}

	fmt.Printf("%v\n", err.Error())

	// Set proper return status or let 500 as default
	returnStatus := http.StatusInternalServerError
	if len(status) > 0 {
		returnStatus = status[0]
	}

	ctrl.SendJSON(w, &msg, returnStatus)
}
