package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/vstdy0/go-shortener/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-shortener/pkg/logging"
)

// newMigrateCmd creates a new migrate command.
func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate DB to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildPsqlStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			ctx, ctxCancel := context.WithTimeout(context.Background(), config.Timeout)
			defer ctxCancel()
			ctx = logging.SetCtxLogger(ctx, logger)

			if err = st.Migrate(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
