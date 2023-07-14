package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

The general form is ` + bold("tucp <resource[ subresource...]> <action> [flags]") + `.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO:
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}
)

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	// rootCmd.AddCommand(
	// 	cmds.WhoAmICmd(),
	// )
}

func main() {
	cobra.EnableCommandSorting = false
	log.SetOutput(os.Stdout)
	rootCmd.Execute()
}
