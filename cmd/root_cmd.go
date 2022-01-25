package cmd

import (
	"github.com/spf13/cobra"
)

var serverConfigFilePath string // only used for serve cmd

var rootCmd = cobra.Command{
	Use: "supabase-admin-api",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	serveCmd.Flags().StringVarP(&serverConfigFilePath, "config-file", "c", "/etc/adminapi/adminapi.yaml", "path to yaml config file")
	rootCmd.AddCommand(&serveCmd, &versionCmd)

	return &rootCmd
}
