package rest

import (
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipDecompressRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func gzipCompressResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	}

	return http.HandlerFunc(fn)
}

type ctxKeyUserID int

const userIDKey ctxKeyUserID = 0
const cookieName = "Authorization"

func cookieAuth(secretKey string) func(next http.Handler) http.Handler {
	cipherKey := sha256.Sum256([]byte(secretKey))
	aesBlock, _ := aes.NewCipher(cipherKey[:])
	aesGCM, _ := cipher.NewGCM(aesBlock)
	nonce := cipherKey[len(cipherKey)-aesGCM.NonceSize():]

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			newCookie := func(w http.ResponseWriter, r *http.Request) {
				userID := uuid.New()
				encryptedValue := aesGCM.Seal(nil, nonce, []byte(userID.String()), nil)

				cookie := http.Cookie{
					Name:  cookieName,
					Value: base64.StdEncoding.EncodeToString(encryptedValue),
					Path:  "/",
				}
				http.SetCookie(w, &cookie)

				ctx = context.WithValue(ctx, userIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				newCookie(w, r)
				return
			}
			encryptedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				newCookie(w, r)
				return
			}
			decryptedValue, err := aesGCM.Open(nil, nonce, encryptedValue, nil)
			if err != nil {
				newCookie(w, r)
				return
			}
			userID, err := uuid.ParseBytes(decryptedValue)
			if err != nil {
				newCookie(w, r)
				return
			}

			ctx = context.WithValue(ctx, userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
