package shortener

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	storagemock "github.com/vstdy0/go-project/storage/mock"
)

type TestSuite struct {
	suite.Suite
	sync.WaitGroup

	svc    *Service
	stMock *storagemock.MockURLStorage

	ctx context.Context
}

func (s *TestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	urlStorageMock := storagemock.NewMockURLStorage(mockCtrl)

	config := Config{
		DelReqTimeout:     5 * time.Second,
		DelBufWipeTimeout: time.Second,
		DelBufCap:         2,
	}

	svc, err := New(
		WithConfig(config),
		WithStorage(urlStorageMock),
	)
	s.Require().NoError(err)

	s.svc = svc
	s.stMock = urlStorageMock
	s.ctx = context.TODO()
}

func TestSuite_URLService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
