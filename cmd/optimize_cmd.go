package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize services",
	Long: `Optimize services for the resources available. 
This should be executed once after any significant resource allocation change.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("optimizeCmd called")
	},
}

func init() {
	rootCmd.AddCommand(optimizeCmd)
}
