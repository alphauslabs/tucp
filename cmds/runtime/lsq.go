package runtime

import (
	"context"
	"io"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

func LsqCmd() *cobra.Command {
	var (
		purge string
	)

	cmd := &cobra.Command{
		Use:   "lsq [next]",
		Short: "Inspect current queue count",
		Long: `Inspect current queue count. Only available in production.

The optional [next] argument is different than the --env=next flag:
the --env=next flag is not supported, while the [next] argument
means query the next environment from prod, which is correct.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			con, err := connection.New(ctx)
			if err != nil {
				logger.Errorf("connection.New failed: %v", err)
				return
			}

			defer con.Close()
			client := tucp.NewTuControlPlaneClient(con)
			req := tucp.ReadProcessQueuesRequest{}
			if len(args) > 0 {
				switch {
				case args[0] != "next":
					logger.Errorf("invalid args: should be 'next'")
					return
				default:
					req.Env = args[0] // next
				}
			}

			req.Purge = purge
			stream, err := client.ReadProcessQueues(ctx, &req)
			if err != nil {
				logger.Errorf("ReadProcessQueues failed: %v", err)
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
				case v.Code != 0:
					logger.Errorf("%v", v.Message)
				default:
					logger.Infof("%v", v.Message)
				}
			}
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&purge, "purge", purge, "comma-separated org list to purge, or 'all'")
	return cmd
}
