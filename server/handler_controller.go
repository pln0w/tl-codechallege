package main

import (
	"net/http"
	"os"
)

// HandlerController structure
type HandlerController struct {
	BaseController
}

type Directory struct {
	Path  string
	Files []string
}

type Watcher struct {
	Addr string
}

// NewHandlerController - returns pointer to HandlerController struct
func NewHandlerController() *HandlerController {
	return &HandlerController{
		BaseController: BaseController{},
	}
}

// Index - function gets request, dispatch proper action to workers via gRPC
// and return results back received content
func (ctrl *HandlerController) Index(w http.ResponseWriter, r *http.Request) {

	var directories []*Directory
	var watchers []*Watcher

	for i := 0; i < len(watchers); i++ {
		directories = append(directories, &Directory{
			Path:  "",
			Files: []string{"", ""},
		})
	}

	// Here call WS breadcast for response

	directories = append(directories, &Directory{
		Path:  "",
		Files: []string{"", ""},
	})

	res := map[string]interface{}{
		"status":   http.StatusOK,
		"response": directories,
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
