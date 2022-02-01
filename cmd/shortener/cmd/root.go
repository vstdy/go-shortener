package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
)

const (
	flagRequestTimeout  = "timeout"
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
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupConfig(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := common.GetConfigFromCmdCtx(cmd)

			opStorage, err := cmd.Flags().GetString(flagStorage)
			if err != nil {
				return fmt.Errorf("app initialization: reading flag %s: %w", flagStorage, err)
			}

			svc, err := cfg.BuildService(opStorage)
			if err != nil {
				return fmt.Errorf("app initialization: service building: %w", err)
			}

			srv := api.NewServer(svc, cfg)

			idleConnsClosed := make(chan struct{})

			go func() {
				sigint := make(chan os.Signal, 1)
				signal.Notify(sigint, os.Interrupt)
				<-sigint

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := srv.Shutdown(ctx); err != nil {
					log.Printf("Server shutdown failed: %v", err)
				}
				log.Println("Server stopped")

				close(idleConnsClosed)
			}()

			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Printf("HTTP server ListenAndServe: %v", err)
			}

			<-idleConnsClosed

			return nil
		},
	}

	cmd.PersistentFlags().Duration(flagRequestTimeout, 5*time.Second, "Request timeout")
	cmd.PersistentFlags().StringP(flagServerAddress, "a", "127.0.0.1:8080", "Server address")
	cmd.PersistentFlags().StringP(flagBaseURL, "b", "http://127.0.0.1:8080", "Base URL")
	cmd.PersistentFlags().StringP(flagFileStoragePath, "f", "./storage/file/storage.txt", "File storage path")
	cmd.PersistentFlags().StringP(flagDatabaseDSN, "d", "", "Database source name")
	cmd.PersistentFlags().StringP(flagStorage, "s", "psql", "Storage type [memory, file, psql]")

	return cmd
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
	config.RequestTimeout = viper.GetDuration(flagRequestTimeout)
	common.SetConfigToCmdCtx(cmd, config)

	return nil
}
