package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
)

const (
	flagLogLevel        = "log-level"
	flagTimeout         = "timeout"
	flagServerAddress   = "server_address"
	flagBaseURL         = "base_url"
	flagFileStoragePath = "file_storage_path"
	envSecretKey        = "secret_key"
	flagDatabaseDSN     = "database_dsn"
	flagStorage         = "storage"
)

// Execute prepares cobra.Command context and executes root cmd.
func Execute() error {
	return newRootCmd().ExecuteContext(common.NewBaseCmdCtx())
}

// newRootCmd creates a new root cmd.
func newRootCmd() *cobra.Command {
	mig := newMigrateCmd()

	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupLogger(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			if err := setupConfig(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)

			opStorage, err := cmd.Flags().GetString(flagStorage)
			if err != nil {
				return fmt.Errorf("app initialization: reading flag %s: %w", flagStorage, err)
			}

			if opStorage == "psql" || opStorage == "" {
				if err = mig.RunE(cmd, args); err != nil {
					return fmt.Errorf("app initialization: migration failed: %w", err)
				}
			}

			svc, err := config.BuildService(opStorage)
			if err != nil {
				return fmt.Errorf("app initialization: service building: %w", err)
			}
			defer func() {
				if err := svc.Close(); err != nil {
					log.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			srv := api.NewServer(svc, config)

			idleConnsClosed := make(chan struct{})

			go func() {
				sigint := make(chan os.Signal, 1)
				signal.Notify(sigint, os.Interrupt)
				<-sigint

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := srv.Shutdown(ctx); err != nil {
					log.Error().Err(err).Msg("Server shutdown failed")
				}
				log.Info().Msg("Server stopped")

				close(idleConnsClosed)
			}()

			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Error().Err(err).Msg("HTTP server ListenAndServe")
			}

			<-idleConnsClosed

			return nil
		},
	}

	cmd.PersistentFlags().String(flagLogLevel, "info", "Logger level [debug,info,warn,error,fatal]")
	cmd.PersistentFlags().Duration(flagTimeout, 5*time.Second, "Request timeout")
	cmd.PersistentFlags().StringP(flagServerAddress, "a", "127.0.0.1:8080", "Server address")
	cmd.PersistentFlags().StringP(flagBaseURL, "b", "http://127.0.0.1:8080", "Base URL")
	cmd.PersistentFlags().StringP(flagFileStoragePath, "f", "./storage/file/storage.txt", "File storage path")
	cmd.PersistentFlags().StringP(flagDatabaseDSN, "d", "", "Database source name")
	cmd.PersistentFlags().StringP(flagStorage, "s", "", "Storage type [memory, file, psql]")

	cmd.AddCommand(mig)

	return cmd
}

// setupLogger configures global logger.
func setupLogger(cmd *cobra.Command) error {
	logLevelBz, err := cmd.Flags().GetString(flagLogLevel)
	if err != nil {
		return fmt.Errorf("%s flag reading: %w", flagLogLevel, err)
	}
	logLevel, err := zerolog.ParseLevel(logLevelBz)
	if err != nil {
		return fmt.Errorf("%s flag parsing: %w", flagLogLevel, err)
	}

	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	log.Logger = log.Output(logWriter).Level(logLevel)

	return nil
}

// setupConfig reads app config and stores it to cobra.Command context.
func setupConfig(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if err := viper.BindEnv(envSecretKey); err != nil {
		return err
	}
	viper.SetDefault(envSecretKey, "secret_key")

	viper.AutomaticEnv()

	config := common.Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("config unmarshal: %w", err)
	}
	config.Timeout = viper.GetDuration(flagTimeout)
	common.SetConfigToCmdCtx(cmd, config)

	return nil
}
