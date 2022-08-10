package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-chi/chi"
)

type LifecycleCommand = string

const LifecycleCommandHeader = "X-Supabase-Lifecycle"

const (
	Stop    LifecycleCommand = "stop"
	Start   LifecycleCommand = "start"
	Restart LifecycleCommand = "restart"
	Enable  LifecycleCommand = "enable"
	Disable LifecycleCommand = "disable"
)

// we default to a restart unless a suitable override is provided
func getLifecycleCommand(r *http.Request) (LifecycleCommand, error) {
	vals, ok := r.Header[LifecycleCommandHeader]
	if !ok || len(vals) == 0 {
		return Restart, nil
	}
	if len(vals) > 1 {
		return Restart, fmt.Errorf("only a single lifecycle command was expected: %+v", vals)
	}
	switch vals[0] {
	case Start:
		return Start, nil
	case Stop:
		return Stop, nil
	case Enable:
		return Enable, nil
	case Disable:
		return Disable, nil
	default:
		return Restart, fmt.Errorf("unknown lifecycle command: %+v", vals[0])
	}
}

// HandleLifecycleCommand is the endpoint for executing service lifecycle commands
func (a *API) HandleLifecycleCommand(w http.ResponseWriter, r *http.Request) error {
	sudo := "sudo"
	app := "systemctl"
	arg0 := "daemon-reload"

	cmd := exec.Command(sudo, app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	fmt.Fprint(os.Stdout, string(stdout))

	lifecycleCommand, err := getLifecycleCommand(r)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return sendJSON(w, http.StatusBadRequest, err.Error())
	}

	// need to do command as goroutine because adminapi gets killed and can't respond
	go func() {
		sudo := "sudo"
		app := "systemctl"
		var arg1 string

		application := chi.URLParam(r, "application")

		switch application {
		case "all":
			arg1 = "services.slice"
		case "gotrue":
			arg1 = "gotrue.service"
		case "postgrest":
			arg1 = "postgrest.service"
		case "pglisten":
			arg1 = "pglisten.service"
		case "kong":
			arg1 = "kong.service"
		case "realtime":
			arg1 = fmt.Sprintf("%s.service", a.config.RealtimeServiceName)
		case "adminapi":
			arg1 = "adminapi.service"
		case "pgsodium":
			arg1 = derivePostgresqlUnitName()
		case "postgresql":
			arg1 = derivePostgresqlUnitName()
		case "pgbouncer":
			arg1 = "pgbouncer.service"
		default:
			arg1 = "services.slice"
		}

		// if admin api is getting restarted give time for http response first
		if application == "adminapi" || application == "all" {
			time.Sleep(2 * time.Second)
		}

		cmd = exec.Command(sudo, app, lifecycleCommand, arg1)
		stdout, err = cmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to %s %s service: %+v\n", lifecycleCommand, arg1, err)
		}

		fmt.Fprint(os.Stdout, string(stdout))
	}()

	return sendJSON(w, http.StatusOK, 200)
}

func derivePostgresqlUnitName() string {
	_, err := os.Stat("/etc/postgresql/postgresql.conf")
	if err != nil {
		return "postgresql@12-main.service"
	} else {
		return "postgresql.service"
	}
}
