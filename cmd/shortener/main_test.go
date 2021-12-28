package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vstdy0/go-project/api"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
	r := api.Router()
	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			method: http.MethodPost,
			path:   "/",
			body:   "https://extremelylengthylink1.com/",
			want: want{
				code:        http.StatusCreated,
				response:    ts.URL + "/1",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			method: http.MethodPost,
			path:   "/",
			body:   "https://extremelylengthylink2.com/",
			want: want{
				code:        http.StatusCreated,
				response:    ts.URL + "/2",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			method: http.MethodGet,
			path:   "/1",
			body:   "",
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.method, tt.path, tt.body)
		defer resp.Body.Close()
		assert.Equal(t, tt.want.code, resp.StatusCode)
		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		assert.Equal(t, tt.want.response, body)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

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
