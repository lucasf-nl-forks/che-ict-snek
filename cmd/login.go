/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [server] [options]",
	Short: "Authenticate with a longboi server",
	Long:  `Authenticate with a longboi server. Supply the root domain of the server you want to connect to, with either the --web or --key option`,
	Run: func(cmd *cobra.Command, args []string) {
		if Web {
			authenticateWithWeb(cmd, args)
		} else if Key != "" {
			authenticateWithKey(cmd, args)
		} else {
			cobra.CheckErr(errors.New("--web or --key option is required"))
		}
	},
}

func authenticateWithWeb(cmd *cobra.Command, args []string) {
	panic("Web authentication not implemented yet, I'm terribly sorry")
}

func authenticateWithKey(cmd *cobra.Command, args []string) {
	viper.Set("ApiServer", args[0])
	viper.Set("ApiKey", Key)
	err := viper.WriteConfig()
	cobra.CheckErr(err)
}

var Key string
var Web bool

func init() {
	authCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringVarP(&Key, "key", "k", "", "Key")
	loginCmd.Flags().BoolVarP(&Web, "web", "w", false, "Web")
}
