package api

import (
	"net/http"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (a *API) PrepareDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/prepare_pg_upgrade.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't unmount data volume")
	}
	logrus.WithField("output", string(output)).Info("Unmounted data volume")

	return nil
}

func (a *API) InitiateDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/initiate_pg_upgrade.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't run pg_upgrade")
	}
	logrus.WithField("output", string(output)).Info("pg_upgrade run finished")

	return nil
}

func (a *API) CompleteDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/complete_pg_upgrade.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't complete pg_upgrade")
	}
	logrus.WithField("output", string(output)).Info("Completed database upgrade")

	return nil
}
