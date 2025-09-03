/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snek",
	Short: "Snek Command Line Tool for Longboi",
	Long:  `Snek is a client for Longboi, together they form a minimal learning platform for programming exercises`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, ".config")
	err = os.MkdirAll(configPath, os.ModeDir)
	cobra.CheckErr(err)

	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("snek")

	err = viper.SafeWriteConfig()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}
}

func init() {

	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose logging")
	cobra.OnInitialize(initConfig)
}
