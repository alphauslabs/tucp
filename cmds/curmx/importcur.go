package curmx

import (
	"context"
	"time"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

func ImportCurCmd() *cobra.Command {
	var (
		force       bool
		noBroadcast bool
	)

	cmd := &cobra.Command{
		Use:   "import <[prefix:]id1[,id2,idn...]> [yyyymm]",
		Short: "Import CUR(s) to our system",
		Long: `Import CUR(s) to our system.

The supported prefix for now is 'orgId:'. Setting 'orgId:{id}' will
import all CURs under that org.

If you provide a list of comma-separated payer ids, make sure they
belong to a single orgId/MSP.

If [yyyymm] is provided, it defaults to the current month.

Set --no-broadcast=true if you want to download the CUR without
triggering our daily calculation batch.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			if len(args) == 0 {
				logger.Errorf("<[prefix:]id1[,id2,idn...]> is required.")
				return
			}

			id := args[0]
			month := time.Now().UTC().Format("200601")
			if len(args) > 1 {
				month = args[1]
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.ImportCurRequest{
				Id:          id,
				Month:       month,
				Force:       force,
				NoBroadcast: noBroadcast,
			}

			resp, err := client.ImportCur(ctx, &req)
			if err != nil {
				logger.Errorf("ImportCur failed: %v", err)
				return
			}

			logger.Infof("Operation [%v] started.", resp.Name)
			logger.Infof("You can check #curimport-trace channel for progress.")
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&force, "force", force, "force import, ignore last modified state")
	cmd.Flags().BoolVar(&noBroadcast, "no-broadcast", noBroadcast, "don't broadcast import event")
	return cmd
}
