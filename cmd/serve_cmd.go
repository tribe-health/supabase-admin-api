package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supabase/supabase-admin-api/api"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	config := api.Config{}
	config.Host = "localhost"
	config.Port = 8085
	config.Endpoint = "ENDPOINT"
	config.RequestIDHeader = "REQUEST_ID_HEADER"
	config.ExternalURL = "API_EXTERNAL_URL"

	api := api.NewAPIWithVersion(&config, Version)

	l := fmt.Sprintf("%v:%v", config.Host, config.Port)
	logrus.Infof("Supabase Admin API started on: %s", l)
	api.ListenAndServe(l)
}
