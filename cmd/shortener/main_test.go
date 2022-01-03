package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/config"
	"github.com/vstdy0/go-project/service/shortener/v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
	var cfg config.Config
	err := env.Parse(&cfg)
	require.NoError(t, err)
	svc, err := shortener.NewService(shortener.WithInMemoryStorage())
	require.NoError(t, err)
	r := api.Router(svc, cfg)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name        string
		method      string
		path        string
		body        string
		contentType string
		want        want
	}{
		{
			method:      http.MethodPost,
			path:        "/",
			body:        "https://extremelylengthylink1.com/",
			contentType: "text/plain; charset=UTF-8",
			want: want{
				code:        http.StatusCreated,
				response:    cfg.BaseURL + "/1",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			method:      http.MethodPost,
			path:        "/api/shorten",
			body:        `{"url": "https://extremelylengthylink2.com/"}`,
			contentType: "application/json",
			want: want{
				code:        http.StatusCreated,
				response:    fmt.Sprintf(`{"result": "%s/%d"}`, cfg.BaseURL, 2),
				contentType: "application/json",
			},
		},
		{
			method:      http.MethodGet,
			path:        "/1",
			body:        "",
			contentType: "",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.method, tt.path, tt.body, tt.contentType)
		defer resp.Body.Close()
		assert.Equal(t, tt.want.code, resp.StatusCode)
		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		assert.Equal(t, tt.want.response, body)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body, contentType string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
