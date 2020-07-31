package cmd

import (
	"context"
	"fmt"

	"github.com/gobuffalo/uuid"
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

func serve(config *Config) {
	ctx, err := api.WithInstanceConfig(context.Background(), config, uuid.Nil)
	if err != nil {
		logrus.Fatalf("Error loading instance config: %+v", err)
	}
	api := api.NewAPIWithVersion(ctx, config, Version)

	l := fmt.Sprintf("%v:%v", config.Host, config.Port)
	logrus.Infof("Supabase Admin API started on: %s", l)
	api.ListenAndServe(l)
}
