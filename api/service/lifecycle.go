package service

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

type LifecycleCommand = string

const (
	Stop    LifecycleCommand = "stop"
	Start   LifecycleCommand = "start"
	Restart LifecycleCommand = "restart"
	Enable  LifecycleCommand = "enable"
	Disable LifecycleCommand = "disable"
)

func ExecuteLifecycleCommand(requestedService string, command LifecycleCommand) error {
	cmd := exec.Command("sudo", "systemctl", "daemon-reload")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return errors.Wrap(err, "failed to reload systemd")
	}
	fmt.Fprint(os.Stdout, string(stdout))

	// need to do command as goroutine because adminapi gets killed and can't respond
	go func() {
		var targetUnit string

		switch requestedService {
		case "all":
			targetUnit = "services.slice"
		case "gotrue":
			targetUnit = "gotrue.service"
		case "postgrest":
			targetUnit = "postgrest.service"
		case "pglisten":
			targetUnit = "pglisten.service"
		case "kong":
			targetUnit = "kong.service"
		case "realtime":
			targetUnit = "realtime.service"
		case "adminapi":
			targetUnit = "adminapi.service"
		case "pgsodium":
			targetUnit = derivePostgresqlUnitName()
		case "postgresql":
			targetUnit = derivePostgresqlUnitName()
		case "pgbouncer":
			targetUnit = "pgbouncer.service"
		default:
			targetUnit = "services.slice"
		}

		// if admin api is getting restarted give time for http response first
		if requestedService == "adminapi" || requestedService == "all" {
			time.Sleep(2 * time.Second)
		}

		cmd = exec.Command("sudo", "systemctl", command, targetUnit)
		stdout, err = cmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to %s %s service: %+v\n", command, targetUnit, err)
		}
		fmt.Fprint(os.Stdout, string(stdout))
	}()
	return nil
}

func derivePostgresqlUnitName() string {
	_, err := os.Stat("/etc/postgresql/postgresql.conf")
	if err != nil {
		return "postgresql@12-main.service"
	} else {
		return "postgresql.service"
	}
}
