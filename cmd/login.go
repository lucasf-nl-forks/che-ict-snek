/*
Copyright © 2025 David Doorn <djdoorn@che.nl>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"snek/utils"
	"sync"
	"time"

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
	srv := &http.Server{Addr: ":9123"}
	var serveRef *http.Server = srv
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		Key = r.URL.Query().Get("key")
		err := utils.ValidateKey(args[0], Key)
		cobra.CheckErr(err)
		viper.Set("ApiServer", args[0])
		viper.Set("ApiKey", Key)
		err = viper.WriteConfig()
		cobra.CheckErr(err)
		log.Println("Logged in successfully")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Logged in successfully, you can close this window now")
		cancel()
	})

	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			cobra.CheckErr(err)
		}
	}()
	dateStr := time.Now().Format("20060102")
	escapedRedirect := url.QueryEscape("http://localhost:9123/token")
	urlStr := fmt.Sprintf("%s/user/apikeys/automated/?name=snek%s&redirect=%s", args[0], dateStr, escapedRedirect)
	fmt.Println("The following url should open in your default browser: ", urlStr)
	err := utils.OpenURL(urlStr)
	cobra.CheckErr(err)
	<-ctx.Done()
	serveRef.Shutdown(ctx)
	wg.Wait()
}

func authenticateWithKey(cmd *cobra.Command, args []string) {

	err := utils.ValidateKey(args[0], Key)
	cobra.CheckErr(err)

	viper.Set("ApiServer", args[0])
	viper.Set("ApiKey", Key)
	err = viper.WriteConfig()
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
