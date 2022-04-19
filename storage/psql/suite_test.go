package psql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/vstdy0/go-shortener/storage/psql/fixtures"
	"github.com/vstdy0/go-shortener/testutils"
)

type TestSuite struct {
	suite.Suite

	container *testutils.PostgreSQLContainer
	storage   *Storage
	fixtures  fixtures.Fixtures

	ctx context.Context
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Minute)
	defer ctxCancel()

	c, err := testutils.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	stCfg := NewDefaultConfig()
	stCfg.DSN = c.GetDSN()

	st, err := NewStorage(WithConfig(stCfg))
	s.Require().NoError(err)

	s.Require().NoError(st.Migrate(ctx))

	fixts, err := fixtures.LoadFixtures(ctx, st.db)
	s.Require().NoError(err)

	s.ctx = context.TODO()
	s.container = c
	s.storage = st
	s.fixtures = fixts
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.container.Terminate(ctx))
}

func TestSuite_PSQLStorage(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
