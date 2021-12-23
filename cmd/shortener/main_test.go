package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		request string
		body    string
		want    want
	}{
		{
			name:    "test #1 | post link",
			method:  http.MethodPost,
			request: "/",
			body:    "https://extremelylengthylink1.com/",
			want: want{
				code:        http.StatusCreated,
				response:    "http://example.com/1",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name:    "test #2 | post link",
			method:  http.MethodPost,
			request: "/",
			body:    "https://extremelylengthylink2.com/",
			want: want{
				code:        http.StatusCreated,
				response:    "http://example.com/2",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name:    "test #3 | get link",
			method:  http.MethodGet,
			request: "/1",
			body:    "",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(shortenURL)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
