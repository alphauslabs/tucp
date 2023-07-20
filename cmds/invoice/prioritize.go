package invoice

import (
	"context"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

func PrioritizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prioritize",
		Short: "Prioritize invoice calculations for a time",
		Long:  `Prioritize invoice calculations over normal daily calculations for a period of time.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			resp, err := client.PrioritizeInvoice(ctx, &tucp.PrioritizeInvoiceRequest{})
			if err != nil {
				logger.Errorf("PrioritizeInvoice failed: %v", err)
				return
			}

			logger.Infof("%v", resp.Status)
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
