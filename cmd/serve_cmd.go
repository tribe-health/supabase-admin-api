package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supabase/supabase-admin-api/api"
	"gopkg.in/yaml.v3"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	bytes, err := os.ReadFile(serverConfigFilePath)
	if err != nil {
		logrus.Fatalf("failed to read in config: %q", err)
	}
	config := api.Config{UpstreamMetricsRefreshDuration: api.DefaultRefreshDuration, RealtimeServiceName: api.DefaultRealtimeServiceName}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		logrus.Fatalf("could not parse config: %q", err)
	}

	createdApiInstance := api.NewAPIWithVersion(&config, Version)
	l := fmt.Sprintf("%v:%v", config.Host, config.Port)
	logrus.Infof("Supabase Admin API %s started on: %s", Version, l)
	createdApiInstance.ListenAndServe(l, config.KeyPath, config.CertPath)
}
