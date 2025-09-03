/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a longboi server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Add credentials to authenticate with a longboi server using `auth login`")
	},
}

func init() {

	rootCmd.AddCommand(authCmd)
}
