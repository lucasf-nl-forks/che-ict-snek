/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"snek/types"
	"snek/utils"
	"time"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Synchronize the current course with the server",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(".snek/snekfile.json")
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No snekfile present, are you in the root of the course directory?")
		}

		var snekFile types.SnekFile
		var updateTime = time.Unix(snekFile.UpdateTime, 0)
		snekFileContent, err := os.ReadFile(".snek/snekfile.json")
		err = json.Unmarshal(snekFileContent, &snekFile)
		if err != nil {
			log.Fatal(err)
		}

		host := viper.GetString("apiServer")
		key := viper.GetString("apiKey")

		content, err := utils.GetCourseContent(host, key, snekFile.CourseSlug)
		if err != nil {
			log.Fatal(err)
		}

		updatedContent := false
		newContent := false
		conflicts := false

		for _, serverExercise := range content {
			found := false
			for _, localExercise := range snekFile.CourseContent {
				if serverExercise.Name == localExercise.Name {
					found = true
					if serverExercise.Hash != localExercise.Hash {
						updatedContent = true
						modTime, _ := getLatestModificationTime(localExercise.Name)
						if modTime.After(updateTime) {
							conflicts = true
							log.Println(fmt.Sprintf("Cannot update exercise %s, since you have changed it. Back up your work and use --force to proceed", localExercise.Name))
						}
					}
				}
			}
			if !found {
				newContent = true
			}
		}

		if !newContent && !updatedContent {
			fmt.Println("You are already up to date :)")
			snekFile.UpdateTime = time.Now().Unix()
			return
		}

		if newContent && (!conflicts || forcePull) {
			for _, serverExercise := range content {
				found := false
				for _, localExercise := range snekFile.CourseContent {
					if serverExercise.Name == localExercise.Name {
						found = true
					}
				}
				if !found {
					getterClient := getter.Client{
						Src:              host + "/" + serverExercise.Url,
						Dst:              serverExercise.Name,
						ProgressListener: utils.NewProgressTracker(),
						Mode:             getter.ClientModeAny,
					}
					err := getterClient.Get()
					if err != nil {
						log.Println("[Error]", err)
					}
				}
			}
		}

		if updatedContent && (!conflicts || forcePull) {
			for _, serverExercise := range content {
				for _, localExercise := range snekFile.CourseContent {
					if serverExercise.Name == localExercise.Name {
						if serverExercise.Hash != localExercise.Hash {
							getterClient := getter.Client{
								Src:              host + "/" + serverExercise.Url,
								Dst:              serverExercise.Name,
								ProgressListener: utils.NewProgressTracker(),
								Mode:             getter.ClientModeAny,
							}
							err := getterClient.Get()
							if err != nil {
								log.Println("[Error]", err)
							}
						}
					}
				}
			}
		}

		if (updatedContent || newContent) && (!conflicts || forcePull) {
			snekFile.CourseContent = content
			snekFile.UpdateTime = time.Now().Unix()
		}

		// finally, write snekfile
		snekFileContent, err = json.MarshalIndent(snekFile, "", "  ")
		_ = os.WriteFile(".snek/snekfile.json", snekFileContent, fs.ModePerm)
	},
}
var forcePull bool

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().BoolVar(&forcePull, "force", false, "force overwrite")
}

func getLatestModificationTime(dirPath string) (time.Time, error) {
	var latestTime time.Time
	err := filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Skipping %q: %v", path, err)
			return nil // Continue walking
		}

		if !entry.IsDir() {
			fileInfo, err := entry.Info()
			if err != nil {
				log.Printf("Skipping file %q: %v", entry.Name(), err)
				return nil
			}

			if latestTime.IsZero() || fileInfo.ModTime().After(latestTime) {
				latestTime = fileInfo.ModTime()
			}
		}

		return nil
	})

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to walk directory: %v", err)
	}

	if latestTime.IsZero() {
		return time.Time{}, fmt.Errorf("no files found")
	}

	return latestTime, nil
}
