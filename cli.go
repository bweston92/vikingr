package main

import (
	"context"
	"errors"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "vikingr [options]",
		Short: "A Viking that swears to protect master branches.",
		Long:  `All master branches for private repos in the given organisation will be protected!`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetString("user") == "" {
				return errors.New("flag user is required")
			}

			if viper.GetString("token") == "" {
				return errors.New("flag token is required")
			}

			if viper.GetString("org") == "" {
				return errors.New("flag org is required")
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

			return runCheck(ctx, gh, org)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("user", "", "", "user to access GitHub as.")
	rootCmd.PersistentFlags().StringP("token", "", "", "token to access GitHub with.")
	rootCmd.PersistentFlags().StringP("org", "", "", "organisation to monitor.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging.")

	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("org", rootCmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

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
