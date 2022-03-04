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
		return errors.Wrap(err, "failed to disable WAL-G")
	}
	logrus.WithField("output", string(output)).Info("WAL-G disabled")
	return nil
}
