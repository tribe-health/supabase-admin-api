package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func DisableRestore(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/disable_walg_restore.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to disable future restorations upon restart")
	}
	logrus.WithField("output", string(output)).Info("WAL-G restore disable completed")
	return nil
}
