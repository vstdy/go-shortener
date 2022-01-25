package root

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func NewRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadEnvs()
			if err != nil {
				return err
			}

			svc, err := shortener.NewService(shortener.WithInFileStorage(cfg))
			if err != nil {
				return err
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

	cfgDefault := config.Config{
		ServerAddress:   "127.0.0.1:8080",
		BaseURL:         "http://127.0.0.1:8080",
		FileStoragePath: "./storage/infile/storage.txt",
		SecretKey:       "secret",
	}

	rootCmd.Flags().StringVarP(&cfgDefault.ServerAddress, "server_address", "a", cfgDefault.ServerAddress, "Set server address")
	rootCmd.Flags().StringVarP(&cfgDefault.BaseURL, "base_url", "b", cfgDefault.BaseURL, "Set base URL")
	rootCmd.Flags().StringVarP(&cfgDefault.FileStoragePath, "file_storage_path", "f", cfgDefault.FileStoragePath, "Set file storage path")
	viper.SetDefault("secret_key", cfgDefault.SecretKey)

	viper.AutomaticEnv()
	if err := viper.BindPFlag("server_address", rootCmd.Flags().Lookup("server_address")); err != nil {
		return nil, err
	}
	if err := viper.BindPFlag("base_url", rootCmd.Flags().Lookup("base_url")); err != nil {
		return nil, err
	}
	if err := viper.BindPFlag("file_storage_path", rootCmd.Flags().Lookup("file_storage_path")); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("secret_key"); err != nil {
		return nil, err
	}

	return rootCmd, nil
}
