package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/cmds/curmx"
	"github.com/alphauslabs/tucp/cmds/invoice"
	"github.com/alphauslabs/tucp/cmds/runtime"
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
		Long: bold("tucp") + ` - Command line interface for our TrueUnblended Control Plane.
Copyright (c) 2023-` + year() + ` Alphaus Cloud, Inc. All rights reserved.

The general form is ` + bold("tucp <resource[ subresource...]> <action> [flags]") + `.

To authenticate, either set GOOGLE_APPLICATION_CREDENTIALS env var or
set the --creds-file flag. You can use the 'iam' tool to request access
to the 'tucpd-[next|prod]' service, like so:

  $ iam allow-me tucpd-prod

You only need to do this once. See https://github.com/alphauslabs/iam
for more information about the 'iam' tool.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			params.ServiceHost = params.HostProd
			switch params.RunEnv {
			case "dev":
				params.ServiceHost = params.HostDev
			case "next":
				params.ServiceHost = params.HostNext
			}

			ctx := context.Background()
			var ts oauth2.TokenSource
			var err error
			switch {
			case params.CredentialsFile != "":
				opts := idtoken.WithCredentialsFile(params.CredentialsFile)
				ts, err = idtoken.NewTokenSource(ctx, "https://"+params.ServiceHost, opts)
			default:
				ts, err = idtoken.NewTokenSource(ctx, "https://"+params.ServiceHost)
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
			if params.Bare {
				logger.SetPrefix(logger.PrefixNone)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().BoolVar(&params.Bare, "bare", params.Bare, "minimal log output")
	rootCmd.PersistentFlags().StringVar(&params.CredentialsFile, "creds-file", "", "optional, GCP service account file")
	rootCmd.PersistentFlags().StringVar(&params.RunEnv, "env", "prod", "dev, next, or prod")
	rootCmd.AddCommand(
		invoice.InvoiceCmd(),
		runtime.RuntimeCmd(),
		curmx.CurmxCmd(),
	)
}

func main() {
	cobra.EnableCommandSorting = false
	log.SetOutput(os.Stdout)
	rootCmd.Execute()
}
