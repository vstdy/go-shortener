package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	servicemock "github.com/vstdy/go-shortener/service/shortener/mock"
)

type TestSuite struct {
	suite.Suite

	srv     *httptest.Server
	svcMock *servicemock.MockService

	config Config
	client *http.Client

	userID uuid.UUID
	cookie *http.Cookie
}

func (s *TestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	ServiceMock := servicemock.NewMockService(mockCtrl)

	config := NewDefaultConfig()
	timeout := time.Second

	clt := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: timeout,
	}

	r, err := NewRouter(ServiceMock, config)
	s.Require().NoError(err)
	srv := httptest.NewServer(r)

	userID := uuid.MustParse("c92d627d-96f1-416f-9c70-75d7504e161e")
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "ov//cWHlDr/8sG2iaSS7/K/QQ/K3X6/cM0cQIXiEZfj0PxLRoh784IF5qDopw1ilxCDMHg==",
		Path:  "/",
	}

	s.srv = srv
	s.svcMock = ServiceMock
	s.config = config
	s.client = clt
	s.userID = userID
	s.cookie = cookie
}

func (s *TestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s TestSuite) testRequest(method, path, body, contentType string) (*http.Response, string) {
	req, err := http.NewRequest(method, s.srv.URL+path, strings.NewReader(body))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(s.cookie)

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestSuite_Server(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
