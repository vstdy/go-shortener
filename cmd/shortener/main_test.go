package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vstdy0/go-project/api"
	"github.com/vstdy0/go-project/cmd/shortener/cmd/common"
)

func TestShortener(t *testing.T) {
	cfg := common.Config{
		RequestTimeout: 5,
		ServerAddress:  "127.0.0.1:8080",
		BaseURL:        "http://127.0.0.1:8080",
		SecretKey:      "test_secret",
	}

	cfg.FileStorage.FileStoragePath = "storage.txt"
	defer os.Remove(cfg.FileStorage.FileStoragePath)

	svc, err := cfg.BuildService("memory")
	require.NoError(t, err)
	r := api.Router(svc, cfg)
	inMemoryTS := httptest.NewServer(r)
	defer inMemoryTS.Close()

	svc, err = cfg.BuildService("file")
	require.NoError(t, err)
	r = api.Router(svc, cfg)
	inFileTS := httptest.NewServer(r)
	defer inFileTS.Close()

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
				response:    fmt.Sprintf(`{"result":"%s/%d"}`, cfg.BaseURL, 2),
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
		memoryResp, memoryBody := testRequest(t, inMemoryTS, tt.method, tt.path, tt.body, tt.contentType)
		defer memoryResp.Body.Close()
		assert.Equal(t, tt.want.code, memoryResp.StatusCode)
		assert.Equal(t, tt.want.contentType, memoryResp.Header.Get("Content-Type"))
		assert.Equal(t, tt.want.response, memoryBody)
		fileResp, fileBody := testRequest(t, inFileTS, tt.method, tt.path, tt.body, tt.contentType)
		defer fileResp.Body.Close()
		assert.Equal(t, tt.want.code, fileResp.StatusCode)
		assert.Equal(t, tt.want.contentType, fileResp.Header.Get("Content-Type"))
		assert.Equal(t, tt.want.response, fileBody)
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
