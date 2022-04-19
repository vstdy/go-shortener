package shortener

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	storagemock "github.com/vstdy0/go-shortener/storage/mock"
)

type TestSuite struct {
	suite.Suite
	sync.WaitGroup

	svc    *Service
	stMock *storagemock.MockStorage

	ctx context.Context
}

func (s *TestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	urlStorageMock := storagemock.NewMockStorage(mockCtrl)

	config := Config{
		DelReqTimeout:     5 * time.Second,
		DelBufWipeTimeout: time.Second,
		DelBufCap:         2,
	}

	svc, err := NewService(
		WithConfig(config),
		WithStorage(urlStorageMock),
	)
	s.Require().NoError(err)

	s.svc = svc
	s.stMock = urlStorageMock
	s.ctx = context.TODO()
}

func TestSuite_Service(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
