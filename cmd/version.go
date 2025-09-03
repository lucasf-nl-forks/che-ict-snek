/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Version string = "unknown"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of snek",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snek version", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
