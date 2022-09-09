package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/supabase/supabase-admin-api/api/service"
)

const LifecycleCommandHeader = "X-Supabase-Lifecycle"

// we default to a restart unless a suitable override is provided
func getLifecycleCommand(r *http.Request) (service.LifecycleCommand, error) {
	vals, ok := r.Header[LifecycleCommandHeader]
	if !ok || len(vals) == 0 {
		return service.Restart, nil
	}
	if len(vals) > 1 {
		return service.Restart, fmt.Errorf("only a single lifecycle command was expected: %+v", vals)
	}
	switch vals[0] {
	case service.Start:
		return service.Start, nil
	case service.Stop:
		return service.Stop, nil
	case service.Enable:
		return service.Enable, nil
	case service.Disable:
		return service.Disable, nil
	default:
		return service.Restart, fmt.Errorf("unknown lifecycle command: %+v", vals[0])
	}
}

// HandleLifecycleCommand is the endpoint for executing service lifecycle commands
func (a *API) HandleLifecycleCommand(w http.ResponseWriter, r *http.Request) error {
	lifecycleCommand, err := getLifecycleCommand(r)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return sendJSON(w, http.StatusBadRequest, err.Error())
	}

	err = service.ExecuteLifecycleCommand(chi.URLParam(r, "application"), lifecycleCommand)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}
	return sendJSON(w, http.StatusOK, 200)
}
