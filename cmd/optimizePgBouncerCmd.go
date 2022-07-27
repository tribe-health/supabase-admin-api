package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/supabase/supabase-admin-api/optimizations"
)

var optimizePgBouncerCmd = &cobra.Command{
	Use:   "pgbouncer",
	Short: "Optimize PGBouncer configuration",
	Long:  `Optimize PGBouncer based on the instance type we're running on.`,
	Run: func(cmd *cobra.Command, args []string) {
		destination := strings.TrimSpace(configFilePath)
		if destination == "" {
			log.Panicln("Destination config file path not specified or invalid.")
		}
		instanceType, err := optimizations.GetInstanceType(context.Background())
		if err != nil {
			log.Panicf("couldn't determine instance type: %+v\n", err)
		}
		err = optimizations.OptimizePgBouncer(destination, *instanceType)
		if err != nil {
			log.Panicln(err)
		}
	},
}

func init() {
	optimizeCmd.AddCommand(optimizePgBouncerCmd)
	optimizePgBouncerCmd.Flags().StringVarP(&configFilePath, "destination-config-file-path", "d", "", "The file we should write the generated configuration to.")
}
