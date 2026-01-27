package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *ServiceTestSuite) SetupSuite() {
	CreatFixtures()
	suite.ctx = context.Background()
}

func (suite *ServiceTestSuite) TearDownSuite() {
	// cleanup
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
