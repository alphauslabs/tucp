package curmx

import (
	"context"
	"io"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

func CurLsqCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsq",
		Short: "Inspect current CUR queue count",
		Long:  `Inspect current CUR queue count. Only available in production.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.ReadCurQueuesRequest{}
			stream, err := client.ReadCurQueues(ctx, &req)
			if err != nil {
				logger.Errorf("ReadCurQueues failed: %v", err)
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
	return cmd
}
