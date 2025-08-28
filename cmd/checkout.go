/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"os"
	"snek/types"
	"snek/utils"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout [course-slug]",
	Short: "Checkout a new course",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Checking out:", args[0])

		var snekFile types.SnekFile
		snekFile.CourseSlug = args[0]

		host := viper.GetString("apiServer")
		key := viper.GetString("apiKey")

		content, err := utils.GetCourseContent(host, key, args[0])
		if err != nil {
			log.Fatal(err)
		}

		snekFile.CourseContent = content

		_ = os.Mkdir(args[0], 0755)
		_ = os.Chdir(args[0])
		_ = os.Mkdir(".snek", 0755)

		for _, contentItem := range content {
			getterClient := getter.Client{
				Src:              host + "/" + contentItem.Url,
				Dst:              contentItem.Name,
				ProgressListener: utils.NewProgressTracker(),
				Mode:             getter.ClientModeAny,
			}
			err := getterClient.Get()
			if err != nil {
				log.Println("[Error]", err)
			}
		}

		snekFile.UpdateTime = time.Now().Unix()

		jsonString, _ := json.MarshalIndent(snekFile, "", "    ")
		_ = os.WriteFile(".snek/snekfile.json", jsonString, 0644)

	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
