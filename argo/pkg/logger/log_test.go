package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SlogSuite struct {
	suite.Suite
}

func TestSlogSuite(t *testing.T) {
	suite.Run(t, new(SlogSuite))
}

func (s *SlogSuite) SetupSuite() {
	Init("info")
}

func (s *SlogSuite) TestInfo() {
	slog.Info("Hello World")
}

func (s *SlogSuite) TestJsonLogger() {
}
