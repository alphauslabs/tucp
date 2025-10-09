package invoice

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alphauslabs/blue-internal-go/tucp/v1"
	"github.com/alphauslabs/bluectl/pkg/logger"
	"github.com/alphauslabs/tucp/pkg/connection"
	"github.com/spf13/cobra"
)

var longRunning bool

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

			if longRunning {
				tick := time.NewTicker(40 * time.Second)
				defer tick.Stop()

				logger.Info("Prioritizing invoice... enter ^C to stop")
				ctxx, cancel := context.WithCancel(ctx)
				defer cancel()

				sigch := make(chan os.Signal, 1)
				go func() {
					signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
					<-sigch
					cancel()
				}()
				_, err := client.PrioritizeInvoice(ctxx, &tucp.PrioritizeInvoiceRequest{})
				if err != nil {
					logger.Errorf("PrioritizeInvoice failed: %v", err)
					return
				}
				for {
					select {
					case <-ctxx.Done():
						logger.Info("Stopping prioritize invoice...")
						return
					case <-tick.C:
						_, err := client.PrioritizeInvoice(ctxx, &tucp.PrioritizeInvoiceRequest{})
						if err != nil {
							logger.Errorf("PrioritizeInvoice failed: %v", err)
							return
						}
					}
				}
			} else {
				resp, err := client.PrioritizeInvoice(ctx, &tucp.PrioritizeInvoiceRequest{})
				if err != nil {
					logger.Errorf("PrioritizeInvoice failed: %v", err)
					return
				}
				logger.Info(resp.Message)
			}
		},
	}
	cmd.Flags().BoolVar(&longRunning, "long", false, "Use this flag when you want to prioritize for a longer period of time")
	cmd.Flags().SortFlags = false
	return cmd
}
