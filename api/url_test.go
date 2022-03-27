package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang/mock/gomock"

	api "github.com/vstdy0/go-shortener/api/model"
	"github.com/vstdy0/go-shortener/model"
	"github.com/vstdy0/go-shortener/pkg"
	serviceMock "github.com/vstdy0/go-shortener/service/shortener/mock"
)

func (s *TestSuite) TestServer_shortenURL() {
	type request struct {
		method      string
		path        string
		body        string
		contentType string
	}

	type expected struct {
		code        int
		prepareBody func(obj model.URL) string
		contentType string
	}

	type testCase struct {
		name         string
		prepareMocks func(ServiceMock *serviceMock.MockService) model.URL
		request      request
		expected     expected
	}
	testCases := []testCase{
		{
			name: "Fail: invalid input",
			prepareMocks: func(ServiceMock *serviceMock.MockService) model.URL {
				input := model.URL{
					UserID: s.userID,
					URL:    "",
				}

				ServiceMock.EXPECT().
					AddURL(gomock.Any(), &input).
					Return(pkg.ErrInvalidInput)

				return model.URL{}
			},
			request: request{
				method:      http.MethodPost,
				path:        "/",
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
			expected: expected{
				code: http.StatusBadRequest,
				prepareBody: func(obj model.URL) string {
					return "invalid input\n"
				},
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Fail: object already exists",
			prepareMocks: func(ServiceMock *serviceMock.MockService) model.URL {
				input := model.URL{
					UserID: s.userID,
					URL:    "https://extremely-lengthy-url.com/",
				}

				ServiceMock.EXPECT().
					AddURL(gomock.Any(), &input).
					Do(func(ctx context.Context, obj *model.URL) {
						obj.ID = 1
					}).
					Return(pkg.ErrAlreadyExists)

				return input
			},
			request: request{
				method:      http.MethodPost,
				path:        "/",
				body:        "https://extremely-lengthy-url.com/",
				contentType: "text/plain; charset=utf-8",
			},
			expected: expected{
				code: http.StatusConflict,
				prepareBody: func(obj model.URL) string {
					obj.ID = 1

					return s.config.BaseURL + "/" + strconv.Itoa(obj.ID)
				},
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "OK",
			prepareMocks: func(ServiceMock *serviceMock.MockService) model.URL {
				input := model.URL{
					UserID: s.userID,
					URL:    "https://extremely-lengthy-url.com/",
				}

				ServiceMock.EXPECT().
					AddURL(gomock.Any(), &input).
					Do(func(ctx context.Context, obj *model.URL) {
						obj.ID = 1
					}).
					Return(nil)

				return input
			},
			request: request{
				method:      http.MethodPost,
				path:        "/",
				body:        "https://extremely-lengthy-url.com/",
				contentType: "text/plain; charset=utf-8",
			},
			expected: expected{
				code: http.StatusCreated,
				prepareBody: func(obj model.URL) string {
					obj.ID = 1

					return s.config.BaseURL + "/" + strconv.Itoa(obj.ID)
				},
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "OK: json request body",
			prepareMocks: func(ServiceMock *serviceMock.MockService) model.URL {
				input := model.URL{
					UserID: s.userID,
					URL:    "https://extremely-lengthy-url.com/",
				}

				ServiceMock.EXPECT().
					AddURL(gomock.Any(), &input).
					Do(func(ctx context.Context, obj *model.URL) {
						obj.ID = 1
					}).
					Return(nil)

				return input
			},
			request: request{
				method:      http.MethodPost,
				path:        "/api/shorten",
				body:        `{"url": "https://extremely-lengthy-url.com/"}`,
				contentType: "application/json",
			},
			expected: expected{
				code: http.StatusCreated,
				prepareBody: func(obj model.URL) string {
					obj.ID = 1

					urlResp := api.NewURLRespFromCanon(obj, s.config.BaseURL)

					resp, err := json.Marshal(urlResp)
					s.Assert().NoError(err)

					return string(resp)
				},
				contentType: "application/json",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.svcMock)

			resp, body := s.testRequest(
				tc.request.method, tc.request.path, tc.request.body, tc.request.contentType)
			defer resp.Body.Close()

			s.Assert().Equal(tc.expected.code, resp.StatusCode)
			s.Assert().Equal(tc.expected.contentType, resp.Header.Get("Content-Type"))
			s.Assert().Equal(tc.expected.prepareBody(input), body)
		})
	}
}

func (s *TestSuite) TestServer_shortenBatchURLs() {
	type request struct {
		method      string
		path        string
		body        string
		contentType string
	}

	type expected struct {
		code        int
		prepareBody func(objs []model.URL) string
		contentType string
	}

	type testCase struct {
		name         string
		prepareMocks func(ServiceMock *serviceMock.MockService) []model.URL
		request      request
		expected     expected
	}
	testCases := []testCase{
		{
			name: "Fail: invalid input",
			prepareMocks: func(ServiceMock *serviceMock.MockService) []model.URL {
				input := []model.URL{
					{
						UserID:        s.userID,
						CorrelationID: "",
						URL:           "https://extremely-lengthy-url.com/",
					},
					{
						UserID:        s.userID,
						CorrelationID: "6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						URL:           "",
					},
				}

				ServiceMock.EXPECT().
					AddBatchURLs(gomock.Any(), &input).
					Return(pkg.ErrInvalidInput)

				return nil
			},
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten/batch",
				body: `
					[
					  {
						"correlation_id":"",
						"original_url":"https://extremely-lengthy-url.com/"
					  },
					  {
						"correlation_id":"6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						"original_url":""
					  }
					]
				`,
				contentType: "application/json",
			},
			expected: expected{
				code: http.StatusBadRequest,
				prepareBody: func(objs []model.URL) string {
					return "invalid input\n"
				},
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Fail: object already exists",
			prepareMocks: func(ServiceMock *serviceMock.MockService) []model.URL {
				input := []model.URL{
					{
						UserID:        s.userID,
						CorrelationID: "056d98a6-f001-4526-b5d9-071900d57363",
						URL:           "https://extremely-lengthy-url-1.com/",
					},
					{
						UserID:        s.userID,
						CorrelationID: "6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						URL:           "https://extremely-lengthy-url-2.com/",
					},
				}

				ServiceMock.EXPECT().
					AddBatchURLs(gomock.Any(), &input).
					Do(func(ctx context.Context, objs *[]model.URL) {
						for idx := range *objs {
							(*objs)[idx].ID = idx + 1
						}
					}).
					Return(pkg.ErrAlreadyExists)

				return input
			},
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten/batch",
				body: `
					[
					  {
						"correlation_id":"056d98a6-f001-4526-b5d9-071900d57363",
						"original_url":"https://extremely-lengthy-url-1.com/"
					  },
					  {
						"correlation_id":"6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						"original_url":"https://extremely-lengthy-url-2.com/"
					  }
					]
				`,
				contentType: "application/json",
			},
			expected: expected{
				code: http.StatusConflict,
				prepareBody: func(objs []model.URL) string {
					for idx := range objs {
						objs[idx].ID = idx + 1
					}

					batchRes := api.NewURLsBatchRespFromCanon(objs, s.config.BaseURL)

					res, err := json.Marshal(batchRes)
					s.Assert().NoError(err)

					return string(res)
				},
				contentType: "application/json",
			},
		},
		{
			name: "OK",
			prepareMocks: func(ServiceMock *serviceMock.MockService) []model.URL {
				input := []model.URL{
					{
						UserID:        s.userID,
						CorrelationID: "056d98a6-f001-4526-b5d9-071900d57363",
						URL:           "https://extremely-lengthy-url-1.com/",
					},
					{
						UserID:        s.userID,
						CorrelationID: "6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						URL:           "https://extremely-lengthy-url-2.com/",
					},
				}

				ServiceMock.EXPECT().
					AddBatchURLs(gomock.Any(), &input).
					Do(func(ctx context.Context, objs *[]model.URL) {
						for idx := range *objs {
							(*objs)[idx].ID = idx + 1
						}
					}).
					Return(nil)

				return input
			},
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten/batch",
				body: `
					[
					  {
						"correlation_id":"056d98a6-f001-4526-b5d9-071900d57363",
						"original_url":"https://extremely-lengthy-url-1.com/"
					  },
					  {
						"correlation_id":"6c9fa3c4-469c-4541-a636-66b7f8b5cbe2",
						"original_url":"https://extremely-lengthy-url-2.com/"
					  }
					]
				`,
				contentType: "application/json",
			},
			expected: expected{
				code: http.StatusCreated,
				prepareBody: func(objs []model.URL) string {
					for idx := range objs {
						objs[idx].ID = idx + 1
					}

					batchRes := api.NewURLsBatchRespFromCanon(objs, s.config.BaseURL)

					res, err := json.Marshal(batchRes)
					s.Assert().NoError(err)

					return string(res)
				},
				contentType: "application/json",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.svcMock)

			resp, body := s.testRequest(
				tc.request.method, tc.request.path, tc.request.body, tc.request.contentType)
			defer resp.Body.Close()

			s.Assert().Equal(tc.expected.code, resp.StatusCode)
			s.Assert().Equal(tc.expected.contentType, resp.Header.Get("Content-Type"))
			s.Assert().Equal(tc.expected.prepareBody(input), body)
		})
	}
}

func (s *TestSuite) TestServer_getShortenedURL() {
	type request struct {
		method string
		path   string
	}

	type expected struct {
		code     int
		location string
	}

	type testCase struct {
		name         string
		prepareMocks func(ServiceMock *serviceMock.MockService)
		request      request
		expected     expected
	}
	testCases := []testCase{
		{
			name: "Fail: invalid input",
			prepareMocks: func(ServiceMock *serviceMock.MockService) {
				input := 1

				ServiceMock.EXPECT().
					GetURL(gomock.Any(), input).
					Return("", pkg.ErrInvalidInput)
			},
			request: request{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: expected{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "OK: no content",
			prepareMocks: func(ServiceMock *serviceMock.MockService) {
				input := 1

				ServiceMock.EXPECT().
					GetURL(gomock.Any(), input).
					Return("", nil)
			},
			request: request{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: expected{
				code:     http.StatusGone,
				location: "",
			},
		},
		{
			name: "OK",
			prepareMocks: func(ServiceMock *serviceMock.MockService) {
				input := 1

				ServiceMock.EXPECT().
					GetURL(gomock.Any(), input).
					Return("https://extremely-lengthy-url.com/", nil)
			},
			request: request{
				method: http.MethodGet,
				path:   "/1",
			},
			expected: expected{
				code:     http.StatusTemporaryRedirect,
				location: "https://extremely-lengthy-url.com/",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepareMocks(s.svcMock)

			resp, _ := s.testRequest(
				tc.request.method, tc.request.path, "", "")
			defer resp.Body.Close()

			s.Assert().Equal(tc.expected.code, resp.StatusCode)
			s.Assert().Equal(tc.expected.location, resp.Header.Get("Location"))
		})
	}
}

func (s *TestSuite) TestServer_getUserURLs() {
	type request struct {
		method string
		path   string
	}

	type expected struct {
		code        int
		prepareBody func(objs []model.URL) string
		contentType string
	}

	type testCase struct {
		name         string
		prepareMocks func(ServiceMock *serviceMock.MockService) []model.URL
		request      request
		expected     expected
	}
	testCases := []testCase{
		{
			name: "OK: no content",
			prepareMocks: func(ServiceMock *serviceMock.MockService) []model.URL {
				ServiceMock.EXPECT().
					GetUserURLs(gomock.Any(), s.userID).
					Return(nil, nil)

				return nil
			},
			request: request{
				method: http.MethodGet,
				path:   "/api/user/urls",
			},
			expected: expected{
				code: http.StatusNoContent,
				prepareBody: func(objs []model.URL) string {
					return ""
				},
				contentType: "",
			},
		},
		{
			name: "OK",
			prepareMocks: func(ServiceMock *serviceMock.MockService) []model.URL {
				output := []model.URL{
					{
						ID:     1,
						UserID: s.userID,
						URL:    "https://extremely-lengthy-url-1.com/",
					},
					{
						ID:     2,
						UserID: s.userID,
						URL:    "https://extremely-lengthy-url-2.com/",
					},
				}

				ServiceMock.EXPECT().
					GetUserURLs(gomock.Any(), s.userID).
					Return(output, nil)

				return output
			},
			request: request{
				method: http.MethodGet,
				path:   "/api/user/urls",
			},
			expected: expected{
				code: http.StatusOK,
				prepareBody: func(objs []model.URL) string {
					userURLs := api.NewUserURLsFromCanon(objs, s.config.BaseURL)

					res, err := json.Marshal(userURLs)
					s.Assert().NoError(err)

					return string(res)
				},
				contentType: "application/json",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			output := tc.prepareMocks(s.svcMock)

			resp, body := s.testRequest(
				tc.request.method, tc.request.path, "", "")
			defer resp.Body.Close()

			s.Assert().Equal(tc.expected.code, resp.StatusCode)
			s.Assert().Equal(tc.expected.contentType, resp.Header.Get("Content-Type"))
			s.Assert().Equal(tc.expected.prepareBody(output), body)
		})
	}
}

func (s *TestSuite) TestServer_deleteUserURLs() {
	type request struct {
		method      string
		path        string
		body        string
		contentType string
	}

	type expected struct {
		code        int
		body        string
		contentType string
	}

	type testCase struct {
		name         string
		prepareMocks func(ServiceMock *serviceMock.MockService)
		request      request
		expected     expected
	}
	testCases := []testCase{
		{
			name: "Fail: invalid input",
			prepareMocks: func(ServiceMock *serviceMock.MockService) {
				ServiceMock.EXPECT().
					RemoveUserURLs(nil).
					Return(pkg.ErrInvalidInput)
			},
			request: request{
				method:      http.MethodDelete,
				path:        "/api/user/urls",
				body:        `[]`,
				contentType: "application/json",
			},
			expected: expected{
				code:        http.StatusBadRequest,
				body:        "invalid input\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "OK",
			prepareMocks: func(ServiceMock *serviceMock.MockService) {
				input := []model.URL{
					{
						ID:     1,
						UserID: s.userID,
					},
					{
						ID:     2,
						UserID: s.userID,
					},
				}

				ServiceMock.EXPECT().
					RemoveUserURLs(input).
					Return(nil)
			},
			request: request{
				method:      http.MethodDelete,
				path:        "/api/user/urls",
				body:        `["1","2"]`,
				contentType: "application/json",
			},
			expected: expected{
				code:        http.StatusAccepted,
				body:        "",
				contentType: "",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepareMocks(s.svcMock)

			resp, body := s.testRequest(
				tc.request.method, tc.request.path, tc.request.body, tc.request.contentType)
			defer resp.Body.Close()

			s.Assert().Equal(tc.expected.code, resp.StatusCode)
			s.Assert().Equal(tc.expected.contentType, resp.Header.Get("Content-Type"))
			s.Assert().Equal(tc.expected.body, body)
		})
	}
}
