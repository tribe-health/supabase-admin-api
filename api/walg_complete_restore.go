package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func CompleteRestorationWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/complete_walg_restore.sh")
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to complete WAL-G restoration"
		logrus.WithError(output).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G restoration complete")
	return nil
}
