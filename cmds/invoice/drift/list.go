package drift

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <orgId:id|compId1[,compIdn]> [yyyymm]",
		Short: "List cost drift information",
		Long: `List cost drift information. If your id has an 'orgId:' prefix, id is
considered an orgId (mspId). If there's no prefix, it is considered
a single companyId or a comma-separated list of companyIds.

If [yyyymm] is not provided, it defaults to the current month.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Error("<id> is required. See -h for usage.")
				return
			}

			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			month := time.Now().UTC().Format("200601")
			if len(args) > 1 {
				month = args[1]
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.ReadCostDriftRequest{Id: args[0], Month: month}
			stream, err := client.ReadCostDrift(ctx, &req)
			if err != nil {
				logger.Errorf("ReadCostDrift failed: %v", err)
				return
			}

		loop:
			for {
				v, err := stream.Recv()
				switch {
				case err == io.EOF:
					break loop
				case err != nil && err != io.EOF:
					logger.Error(err)
					break loop
				}

				b, _ := json.Marshal(v)
				logger.Info(string(b))
			}
		},
	}

	cmd.Flags().SortFlags = false
	return cmd
}
