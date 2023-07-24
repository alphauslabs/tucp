package curmx

import (
	"context"
	"io"
	"time"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

func SimulateCurImportedCmd() *cobra.Command {
	var (
		overrideDates string
	)

	cmd := &cobra.Command{
		Use:   "simulate-imported <[prefix:]id> [yyyymm]",
		Short: "Simulate a CUR imported event",
		Long: `Simulate a CUR imported event. Useful if you want to trigger a payer-level
daily calculation without actually redownloading the CUR.

The supported prefixes are 'org:' and 'payer:'. No prefix defaults
to payer. If [yyyymm] is not set, it defaults to the current month.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			if len(args) == 0 {
				logger.Errorf("<[prefix:]id> is required.")
				return
			}

			id := args[0]
			month := time.Now().UTC().Format("200601")
			if len(args) > 1 {
				month = args[1]
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.SimulateCurImportedRequest{Id: id, Month: month}
			req.OverrideDates = overrideDates
			stream, err := client.SimulateCurImported(ctx, &req)
			if err != nil {
				logger.Errorf("SimulateCurImportedRequest failed: %v", err)
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

				switch {
				case v.Message.Code != int32(codes.OK):
					logger.Error(v.Message.Message)
				default:
					logger.Info(v.Message.Message)
				}
			}
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&overrideDates, "override-dates", "", "comma-separated dates to simulate (yyyymmdd)")
	return cmd
}
