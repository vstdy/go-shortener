package grpc

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	HeaderAuthorize = "authorization"
)

type ctxKeyUserID int

const userIDKey ctxKeyUserID = 0

func metadataAuth(secretKey string) func(ctx context.Context) (context.Context, error) {
	cipherKey := sha256.Sum256([]byte(secretKey))
	aesBlock, _ := aes.NewCipher(cipherKey[:])
	aesGCM, _ := cipher.NewGCM(aesBlock)
	nonce := cipherKey[len(cipherKey)-aesGCM.NonceSize():]

	return func(ctx context.Context) (context.Context, error) {
		newCreds := func(ctx context.Context) (context.Context, error) {
			userID := uuid.New()
			encryptedValue := aesGCM.Seal(nil, nonce, []byte(userID.String()), nil)
			token := base64.StdEncoding.EncodeToString(encryptedValue)

			header := metadata.New(map[string]string{HeaderAuthorize: token})
			grpc.SetHeader(ctx, header)

			ctx = metadata.AppendToOutgoingContext(ctx, HeaderAuthorize, token)
			ctx = context.WithValue(ctx, userIDKey, userID)

			return ctx, nil
		}

		token := metautils.ExtractIncoming(ctx).Get(HeaderAuthorize)
		if token == "" {
			return newCreds(ctx)
		}
		encryptedValue, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return newCreds(ctx)
		}
		decryptedValue, err := aesGCM.Open(nil, nonce, encryptedValue, nil)
		if err != nil {
			return newCreds(ctx)
		}
		userID, err := uuid.ParseBytes(decryptedValue)
		if err != nil {
			return newCreds(ctx)
		}

		header := metadata.New(map[string]string{HeaderAuthorize: token})
		grpc.SetHeader(ctx, header)

		ctx = context.WithValue(ctx, userIDKey, userID)

		return ctx, nil
	}
}

func timeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		_ interface{}, err error) {

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				err = status.Error(codes.DeadlineExceeded, "context deadline exceeded")
			}
		}()

		return handler(ctx, req)
	}
}
