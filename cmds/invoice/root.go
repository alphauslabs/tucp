package invoice

import (
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/cmds/invoice/drift"
	"github.com/spf13/cobra"
)

func InvoiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "Subcommand for invoice-related operations",
		Long:  `Subcommand for invoice-related operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(
		StartCmd(),
		PrioritizeCmd(),
		drift.DriftCmd(),
	)

	return cmd
}
