package api

import (
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func cookieAuth(userID int, secretKey string) func(next http.Handler) http.Handler {
	cookieName := "Authentication"
	cipherKey := sha256.Sum256([]byte(secretKey))
	aesBlock, _ := aes.NewCipher(cipherKey[:])
	aesGCM, _ := cipher.NewGCM(aesBlock)
	nonce := cipherKey[len(cipherKey)-aesGCM.NonceSize():]

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			switch cookie, _ := r.Cookie(cookieName); {
			case cookie != nil:
				encryptedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
				if err != nil {
					break
				}
				decryptedValue, err := aesGCM.Open(nil, nonce, encryptedValue, nil)
				if err != nil {
					break
				}

				ctx = context.WithValue(ctx, userIDKey, string(decryptedValue))
				next.ServeHTTP(w, r.WithContext(ctx))
			default:
				userID++
				encryptedValue := aesGCM.Seal(nil, nonce, []byte(strconv.Itoa(userID)), nil)

				cookie := http.Cookie{
					Name:  cookieName,
					Value: base64.StdEncoding.EncodeToString(encryptedValue),
				}
				http.SetCookie(w, &cookie)

				ctx = context.WithValue(ctx, userIDKey, strconv.Itoa(userID))
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		return http.HandlerFunc(fn)
	}
}
