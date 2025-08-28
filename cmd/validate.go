/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"fmt"
	"snek/utils"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check whether you're logged in",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("apiServer")
		key := viper.GetString("apiKey")
		if host == "" || key == "" {
			cobra.CheckErr(fmt.Errorf("Login is not configured"))
			return
		}

		err := utils.ValidateKey(host, key)
		cobra.CheckErr(err)
		fmt.Println("Login OK")
	},
}

func init() {
	authCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
