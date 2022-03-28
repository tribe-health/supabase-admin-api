package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func DisableWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/disable_walg.sh")
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to disable WAL-G"
		logrus.WithError(output).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G disabled")
	return nil
}
