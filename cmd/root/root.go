package root

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vstdy0/go-project/config"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.AutomaticEnv()
			err := viper.BindPFlag("server_address", cmd.Flags().Lookup("server_address"))
			if err != nil {
				return err
			}
			err = viper.BindPFlag("base_url", cmd.Flags().Lookup("base_url"))
			if err != nil {
				return err
			}
			err = viper.BindPFlag("file_storage_path", cmd.Flags().Lookup("file_storage_path"))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cfgDefault := config.Config{
		ServerAddress:   "127.0.0.1:8080",
		BaseURL:         "http://127.0.0.1:8080",
		FileStoragePath: "./storage/infile/storage.txt",
	}
	flags := rootCmd.Flags()
	flags.StringVarP(&cfgDefault.ServerAddress, "server_address", "a", cfgDefault.ServerAddress, "Set server address")
	flags.StringVarP(&cfgDefault.BaseURL, "base_url", "b", cfgDefault.BaseURL, "Set base URL")
	flags.StringVarP(&cfgDefault.FileStoragePath, "file_storage_path", "f", cfgDefault.FileStoragePath, "Set file storage path")

	return rootCmd
}
