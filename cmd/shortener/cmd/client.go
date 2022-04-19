package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	apiGrpc "github.com/vstdy0/go-shortener/api/grpc"
	"github.com/vstdy0/go-shortener/api/grpc/model"
	"github.com/vstdy0/go-shortener/cmd/shortener/cmd/common"
	"github.com/vstdy0/go-shortener/pkg/grpc/url-service"
	"github.com/vstdy0/go-shortener/pkg/logging"
)

const (
	flagToken = "token"
)

// newClientCmd creates a new gRPC-client command.
func newClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Debug gRPC client",
	}

	cmd.PersistentFlags().StringP(flagToken, "t", "", "Authentication token")

	cmd.AddCommand(shortenURLCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(deleteUserURLsCmd())

	return cmd
}

// newClientCmd creates a gRPC-client command for get requests.
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Client get commands",
	}

	cmd.AddCommand(getOriginalURLCmd())
	cmd.AddCommand(getUsersURLsCmd())

	return cmd
}

// shortenURLCmd returns a gRPC-client command for ShortenURL request.
func shortenURLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shorten",
		Short:   "Shorten given URL",
		Example: "shorten {url}",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))
			ctx := logging.SetCtxLogger(context.Background(), logger)

			conn, err := createGRPCClientConnection(config.GRPCServer.ServerAddress, logger)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := urlService.NewURLServiceClient(conn)

			ctx, err = parseTokenFlag(cmd, ctx)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			var header metadata.MD
			resp, err := client.ShortenURL(
				ctx,
				model.NewShortenURLReq(args[0]),
				grpc.Header(&header),
			)
			if err != nil {
				return fmt.Errorf("request failed: %v", err)
			}

			logger.Info().Msgf("%s\ntoken %s", resp.String(), header[apiGrpc.HeaderAuthorize][0])

			return nil
		},
	}

	cmd.AddCommand(shortenURLsBatchCmd())

	return cmd
}

// shortenURLsBatchCmd returns a gRPC-client command for ShortenURLsBatch request.
func shortenURLsBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "batch",
		Short:   "Shorten given URLs batch",
		Example: "shorten batch {url}...",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))
			ctx := logging.SetCtxLogger(context.Background(), logger)

			conn, err := createGRPCClientConnection(config.GRPCServer.ServerAddress, logger)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := urlService.NewURLServiceClient(conn)

			ctx, err = parseTokenFlag(cmd, ctx)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			var header metadata.MD
			resp, err := client.ShortenURLsBatch(
				ctx,
				model.NewShortenURLsBatchReq(args),
				grpc.Header(&header),
			)
			if err != nil {
				return fmt.Errorf("request failed: %v", err)
			}

			logger.Info().Msgf("%s\ntoken %s", resp.String(), header[apiGrpc.HeaderAuthorize][0])

			return nil
		},
	}

	return cmd
}

// getOriginalURLCmd returns a gRPC-client command for GetOriginalURL request.
func getOriginalURLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "original_url",
		Short:   "Get original URL",
		Example: "get original_url {url}",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))
			ctx := logging.SetCtxLogger(context.Background(), logger)

			conn, err := createGRPCClientConnection(config.GRPCServer.ServerAddress, logger)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := urlService.NewURLServiceClient(conn)

			ctx, err = parseTokenFlag(cmd, ctx)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			var header metadata.MD
			resp, err := client.GetOriginalURL(
				ctx,
				model.NewGetOrigURLReq(args[0], config.GRPCServer.BaseURL),
				grpc.Header(&header),
			)
			if err != nil {
				return fmt.Errorf("request failed: %v", err)
			}

			logger.Info().Msgf("%s\ntoken %s\n%s", resp.String(),
				header[apiGrpc.HeaderAuthorize][0], header["location"][0],
			)

			return nil
		},
	}

	return cmd
}

// getUsersURLsCmd returns a gRPC-client command for GetUsersURLs request.
func getUsersURLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users_urls",
		Short:   "Get user's URLs",
		Example: "get users_urls",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))
			ctx := logging.SetCtxLogger(context.Background(), logger)

			conn, err := createGRPCClientConnection(config.GRPCServer.ServerAddress, logger)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := urlService.NewURLServiceClient(conn)

			ctx, err = parseTokenFlag(cmd, ctx)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			var header metadata.MD
			resp, err := client.GetUsersURLs(
				ctx,
				&emptypb.Empty{},
				grpc.Header(&header),
			)
			if err != nil {
				return fmt.Errorf("request failed: %v", err)
			}

			logger.Info().Msgf("%s\ntoken %s", resp.GetResponse(), header[apiGrpc.HeaderAuthorize][0])

			return nil
		},
	}

	return cmd
}

// deleteUserURLsCmd returns a gRPC-client command for DeleteUserURLs request.
func deleteUserURLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete given users URLs",
		Example: "delete [ids]...",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))
			ctx := logging.SetCtxLogger(context.Background(), logger)

			conn, err := createGRPCClientConnection(config.GRPCServer.ServerAddress, logger)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := urlService.NewURLServiceClient(conn)

			ctx, err = parseTokenFlag(cmd, ctx)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			var header metadata.MD
			_, err = client.DeleteUserURLs(
				ctx,
				model.NewDelUserURLsReq(args),
				grpc.Header(&header),
			)
			if err != nil {
				return fmt.Errorf("request failed: %v", err)
			}

			logger.Info().Msgf("%s\ntoken %s", "ok", header[apiGrpc.HeaderAuthorize][0])

			return nil
		},
	}

	return cmd
}

// parseTokenFlag adds authentication to header.
func parseTokenFlag(cmd *cobra.Command, ctx context.Context) (context.Context, error) {
	token, err := cmd.Flags().GetString(flagToken)
	if err != nil {
		return nil, fmt.Errorf("parsing '%s' flag: %v", flagToken, err)
	}

	return metadata.AppendToOutgoingContext(ctx, apiGrpc.HeaderAuthorize, token), nil
}

// createGRPCClientConnection creates a new gRPC client connection.
func createGRPCClientConnection(serverAddr string, logger zerolog.Logger) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("creating client connection failed: %v", err)
	}

	logger.Info().Msgf("client gRPC: %s", serverAddr)

	return conn, nil
}
