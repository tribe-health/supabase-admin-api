package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func EnableWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/enable_walg.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to enable WAL-G")
	}
	logrus.WithField("output", string(output)).Info("WAL-G enabled")
	return nil
}
