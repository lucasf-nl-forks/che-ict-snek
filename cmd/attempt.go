/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"snek/utils"
	"strings"
	"time"
)

// attemptCmd represents the attempt command
var attemptCmd = &cobra.Command{
	Use:   "attempt",
	Short: "Attempt the Exercise in the current directory, pass -e or --exercise if the Exercise cannot be determined",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := os.Getwd()
		if Exercise == "" {
			if err != nil {
				return
			}
			// Get the last two parts of the working directory as identifier
			workDir = strings.ReplaceAll(workDir, "\\", "/")
			workDirParts := strings.Split(workDir, "/")
			Exercise = workDirParts[len(workDirParts)-2] + "/" + workDirParts[len(workDirParts)-1]
		}
		log.Println("Found identifier:", Exercise)
		apiServer := viper.GetString("apiServer")
		apiKey := viper.GetString("apiKey")
		endpoint := apiServer + "/api/attempt/exercise/" + Exercise
		spin := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
		spin.Suffix = "\t Submitting Attempt"
		spin.Start()
		attemptData, submitErr := utils.ZipAndSubmitAttempt(workDir, endpoint, apiKey)
		spin.Stop()
		cobra.CheckErr(submitErr)
		log.Println("Attempt id:", attemptData.Slug)
		spin.UpdateCharSet(spinner.CharSets[11])
		spin.Suffix = "\t Waiting for attempt to run"
		spin.Start()
		var attemptFinished = false
		var requestError error
		for !attemptFinished {
			time.Sleep(5 * time.Second)
			attemptData, requestError = utils.CheckAttemptStatus(apiServer+"/api/attempt/"+attemptData.Slug, apiKey)
			cobra.CheckErr(requestError)
			if attemptData.Status == "finished" {
				attemptFinished = true
			}
			if attemptData.Status == "running" {
				spin.Stop()
				spin.Suffix = "\t Running Attempt"
				spin.Start()
			}
		}
		spin.Stop()
		fmt.Printf("Finished running in %d seconds \r\n", attemptData.Runtime)
		if attemptData.Succeeded {
			fmt.Println("✅ Attempt succeeded")
		} else {
			fmt.Println("❌ Attempt failed")
			fmt.Println(attemptData.Output)
		}
	},
}

var Exercise string

func init() {
	rootCmd.AddCommand(attemptCmd)

	attemptCmd.Flags().StringVarP(&Exercise, "exercise", "e", "", "Exercise to run")
}
