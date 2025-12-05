package config

import (
	"log/slog"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoaderSuite struct {
	suite.Suite
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}

func (s *LoaderSuite) SetupSuite() {
	err := Init("../../configs", "config", "ARGO")
	if err != nil {
		panic(err)
	}
}

func (s *LoaderSuite) TestLoadConfig() {
	type Config struct {
		DB struct {
			Host     string `yaml:"host"`
			Database string `yaml:"database"`
		}
		Server struct {
			Port string `yaml:"port"`
		}
	}
	var cf Config
	err := viper.Unmarshal(&cf)
	if err != nil {
		slog.Error("unable to decode into struct", "err", err)
	}
	assert.Equal(s.T(), cf.DB.Host, "localhost")
	assert.Equal(s.T(), cf.DB.Database, "test1")
	assert.Equal(s.T(), cf.Server.Port, "9090")
}

func (s *LoaderSuite) TestEnv() {
	assert.Equal(s.T(), viper.GetString("server.port"), "9090")
	assert.Equal(s.T(), viper.GetString("db.mode"), "test1")
}
