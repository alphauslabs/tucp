package invoice

import (
	"context"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	var (
		compIds     string
		force       bool
		acctsOnly   bool
		tagsOnly    bool
		skipFees    bool
		skipSupport bool
	)

	cmd := &cobra.Command{
		Use:   "start <orgId> [yyyymm]",
		Short: "Start an invoice calculation",
		Long: `Start an invoice calculation. If [yyyymm] is not set, calculate
for the previous UTC month.

If both --accts-only and --tags-only are set, --tags-only is
discarded and will proceed calculations as --accts-only.

The --skip-support functions include (but not limited to) raw
unblended exports, RI+SP and invoice id detections.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Error("[orgId] is required. See -h for usage.")
				return
			}

			orgId := args[0]
			var month string
			if len(args) > 1 {
				month = args[1]
			}

			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.StartInvoiceRequest{
				OrgId:        orgId,
				Month:        month,
				Force:        force,
				AccountsOnly: acctsOnly,
				TagsOnly:     tagsOnly,
				SkipFees:     skipFees,
				SkipSupport:  skipSupport,
			}

			resp, err := client.StartInvoice(ctx, &req)
			if err != nil {
				logger.Errorf("StartInvoice failed: %v", err)
				return
			}

			logger.Infof("Operation [%v] started.", resp.Name)
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&compIds, "companies", compIds, "comma-separated list of company ids")
	cmd.Flags().BoolVar(&force, "force", force, "overwrite existing calculation, if any")
	cmd.Flags().BoolVar(&acctsOnly, "accts-only", acctsOnly, "only for account BGs, skip tags")
	cmd.Flags().BoolVar(&tagsOnly, "tags-only", tagsOnly, "only for tags-based BGs, skip accounts")
	cmd.Flags().BoolVar(&skipFees, "skip-fees", skipFees, "skip fees aggregation/calculation")
	cmd.Flags().BoolVar(&skipSupport, "skip-support", skipSupport, "skip additional support functions")
	return cmd
}
