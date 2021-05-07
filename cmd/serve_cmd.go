package cmd

import (
	"fmt"
	"os"
	"strconv"

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

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func serve() {
	config := api.Config{}
	config.Host = getEnv("HOST", "localhost")
	port, _ := strconv.Atoi(getEnv("PORT", "8085"))
	config.Port = port
	config.Endpoint = "ENDPOINT"
	config.RequestIDHeader = "REQUEST_ID_HEADER"
	config.ExternalURL = "API_EXTERNAL_URL"

	api := api.NewAPIWithVersion(&config, Version)

	l := fmt.Sprintf("%v:%v", config.Host, config.Port)
	logrus.Infof("Supabase Admin API started on: %s", l)
	api.ListenAndServe(l)
}
