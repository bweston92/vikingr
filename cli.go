package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vikingr [options] [command]",
		Short: "A Viking that swears to protect master branches.",
		Long:  `All master branches for private repos in the given organisation will be protected!`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("user") == "" {
				return errors.New("user is required")
			}

			if viper.GetString("token") == "" {
				return errors.New("token is required")
			}

			if viper.GetString("org") == "" {
				return errors.New("org is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			gh := getGithubClient(
				viper.GetString("user"),
				viper.GetString("token"),
			)
			org := viper.GetString("org")

			frequency := viper.GetInt64("frequency")
			if frequency == 0 {
				return runCheck(ctx, gh, org)
			}

			ticker := time.NewTicker(time.Duration(frequency) * time.Minute)
			quit := make(chan os.Signal, 1)
			done := make(chan bool, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				for {
					select {
					case <-ticker.C:
						if err := runCheck(ctx, gh, org); err != nil {
							log.Printf("error(s) occured during protecting phase: %s\n", err)
						}
					case <-quit:
						ticker.Stop()
						done <- true
						return
					}
				}
			}()

			<-done
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("user", "", "", "user to access GitHub as.")
	rootCmd.PersistentFlags().StringP("token", "", "", "token to access GitHub with.")
	rootCmd.PersistentFlags().StringP("org", "", "", "organisation to monitor.")
	rootCmd.PersistentFlags().Int64P("frequency", "", 0, "frequency to run in minutes")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging.")

	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("org", rootCmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("frequency", rootCmd.PersistentFlags().Lookup("frequency"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	viper.SetEnvPrefix("VIKINGR")
	viper.AutomaticEnv()
}

func debug(f string, v ...interface{}) {
	if viper.GetBool("verbose") {
		if len(v) > 0 {
			log.Printf(f+"\n", v...)
		} else {
			log.Println(f)
		}
	}
}
