package cmd

import (
	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize services",
	Long: `Optimize services for the resources available. 
This should be executed once after any significant resource allocation change.`,
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
}
