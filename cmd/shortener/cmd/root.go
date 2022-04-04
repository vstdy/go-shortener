package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy0/go-shortener/api/rest"
	"github.com/vstdy0/go-shortener/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-shortener/pkg/logging"
	"github.com/vstdy0/go-shortener/pkg/tracing"
)

const (
	flagConfigPath      = "config"
	flagLogLevel        = "log_level"
	flagTimeout         = "timeout"
	flagServerAddress   = "server_address"
	flagBaseURL         = "base_url"
	flagFileStoragePath = "file_storage_path"
	envSecretKey        = "secret_key"
	flagDatabaseDSN     = "database_dsn"
	flagStorageType     = "storage_type"
)

// Execute prepares cobra.Command context and executes root cmd.
func Execute() error {
	return newRootCmd().ExecuteContext(common.NewBaseCmdCtx())
}

// newRootCmd creates a new root cmd.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupConfig(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))

			tracerCloser, err := tracing.SetupGlobalJaegerTracer()
			if err != nil {
				return fmt.Errorf("app initialization: tracer setting: %w", err)
			}
			defer tracerCloser.Close()

			svc, err := config.BuildService()
			if err != nil {
				return fmt.Errorf("app initialization: service building: %w", err)
			}

			srv, err := rest.NewServer(svc, config.Server)
			if err != nil {
				return fmt.Errorf("app initialization: server building: %w", err)
			}

			go func() {
				if err = srv.ListenAndServe(); err != http.ErrServerClosed {
					logger.Error().Err(err).Msg("HTTP server ListenAndServe")
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)
			<-stop

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err = srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}
			if err = svc.Close(); err != nil {
				return fmt.Errorf("service shutdown failed: %w", err)
			}
			logger.Info().Msg("server stopped")

			return nil
		},
	}

	config := common.BuildDefaultConfig()
	cmd.PersistentFlags().String(flagConfigPath, "./config.toml", "Config file path")
	cmd.PersistentFlags().StringP(flagLogLevel, "l", config.LogLevel.String(), "Logger level [debug,info,warn,error,fatal]")
	cmd.PersistentFlags().Duration(flagTimeout, config.Timeout, "Request timeout")
	cmd.PersistentFlags().StringP(flagDatabaseDSN, "d", config.PSQLStorage.DSN, "Database source name")
	cmd.Flags().StringP(flagServerAddress, "a", config.Server.ServerAddress, "Server address")
	cmd.Flags().StringP(flagBaseURL, "b", config.Server.BaseURL, "Base URL")
	cmd.Flags().StringP(flagStorageType, "s", config.StorageType, "Storage type [memory, file, psql]")
	cmd.Flags().StringP(flagFileStoragePath, "f", config.FileStorage.FileStoragePath, "File storage path")

	cmd.AddCommand(newMigrateCmd())

	return cmd
}

// setupConfig reads app config and stores it to cobra.Command context.
func setupConfig(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("flags binding: %w", err)
	}

	if err := viper.BindEnv(envSecretKey); err != nil {
		return fmt.Errorf("%s env binding: %w", envSecretKey, err)
	}

	configPath := viper.GetString(flagConfigPath)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	config := common.BuildDefaultConfig()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("config unmarshal: %w", err)
	}

	config.Timeout = viper.GetDuration(flagTimeout)
	config.Server.Timeout = config.Timeout

	logLevel, err := zerolog.ParseLevel(viper.GetString(flagLogLevel))
	if err != nil {
		return fmt.Errorf("%s flag parsing: %w", flagLogLevel, err)
	}
	config.LogLevel = logLevel
	config.Server.LogLevel = config.LogLevel

	common.SetConfigToCmdCtx(cmd, config)

	return nil
}
