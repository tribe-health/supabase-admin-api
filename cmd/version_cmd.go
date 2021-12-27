package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
)

//go:embed VERSION
var Version string

var versionCmd = cobra.Command{
	Run: showVersion,
	Use: "version",
}

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}
