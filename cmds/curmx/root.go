package curmx

import (
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/spf13/cobra"
)

func CurmxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "curmx",
		Short: "Subcommand for CUR import/export operations",
		Long:  `Subcommand for CUR import/export operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("see -h for more information")
		},
	}

	cmd.Flags().SortFlags = false
	cmd.AddCommand(SimulateCurImportedCmd())
	return cmd
}
