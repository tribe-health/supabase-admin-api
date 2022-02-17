package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func RestoreDatabase(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/commence_walg_restore.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to execute WAL-G restoration")
	}
	logrus.WithField("output", string(output)).Info("WAL-G restoration completed")
	return nil
}
