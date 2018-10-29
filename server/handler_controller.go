package main

import (
	"net/http"
	"os"
)

// HandlerController structure
type HandlerController struct {
	BaseController
}

// Directory structure
type Directory struct {
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

// NewHandlerController - returns pointer to HandlerController struct
func NewHandlerController() *HandlerController {
	return &HandlerController{
		BaseController: BaseController{},
	}
}

// Index - function gets request, dispatch proper action to workers
// and return results back received content
func (ctrl *HandlerController) Index(w http.ResponseWriter, r *http.Request) {

	var directories []*Directory

	for c := range hub.clients {
		directories = append(directories, &Directory{
			Path:  c.dir,
			Files: c.files,
		})
	}

	res := map[string]interface{}{
		"status":      http.StatusOK,
		"directories": directories,
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}

// HealthCheck - function returns status and hostname
func (ctrl *HandlerController) HealthCheck(w http.ResponseWriter, r *http.Request) {

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

// DumpWatchers - function returns connected watchers details
func (ctrl *HandlerController) DumpWatchers(w http.ResponseWriter, r *http.Request) {

	res := map[string]interface{}{
		"status":   http.StatusOK,
		"watchers": hub.getClients(),
	}

	ctrl.SendJSON(w, &res, http.StatusOK)
}
