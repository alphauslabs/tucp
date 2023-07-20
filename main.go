package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/params"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

var (
	bold = color.New(color.Bold).SprintFunc()
	year = func() string {
		return fmt.Sprintf("%v", time.Now().Year())
	}

	rootCmd = &cobra.Command{
		Use:   "tucp",
		Short: bold("tucp") + " - Command line interface for tucpd",
		Long: bold("tucp") + ` - Command line interface for the TrueUnblended Control Plane.
Copyright (c) 2023-` + year() + ` Alphaus Cloud, Inc. All rights reserved.

The general form is ` + bold("tucp <resource[ subresource...]> <action> [flags]") + `.

To authenticate, either set GOOGLE_APPLICATION_CREDENTIALS env var or
set the --svcacct-file flag. Ask the service owner if your credentials
file doesn't have access to the service itself.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// dev: tucpd-dev-cnugyv5cta-an.a.run.app
			// next: tucpd-next-u554nqhjka-an.a.run.app
			// prod: tucpd-prod-u554nqhjka-an.a.run.app
			svc := "tucpd-prod-u554nqhjka-an.a.run.app"
			ctx := context.Background()
			var ts oauth2.TokenSource
			var err error

			switch {
			case params.CredentialsFile != "":
				opts := idtoken.WithCredentialsFile(params.CredentialsFile)
				ts, err = idtoken.NewTokenSource(ctx, "https://"+svc, opts)
			default:
				ts, err = idtoken.NewTokenSource(ctx, "https://"+svc)
			}

			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			token, err := ts.Token()
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}

			params.AccessToken = token.AccessToken
			logger.Info(params.AccessToken)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().StringVar(&params.CredentialsFile, "svcacct-file", "", "GCP service account file")
	rootCmd.PersistentFlags().StringVar(&params.AccessToken, "access-token", "", "use directly if not empty")
	// rootCmd.AddCommand(
	// 	cmds.WhoAmICmd(),
	// )
}

func main() {
	cobra.EnableCommandSorting = false
	log.SetOutput(os.Stdout)
	rootCmd.Execute()
}
