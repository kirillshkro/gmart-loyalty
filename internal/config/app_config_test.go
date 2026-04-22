package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AppConfigSuite struct {
	suite.Suite
}

func (a *AppConfigSuite) SetupSuite() {

}

func (a *AppConfigSuite) TearDownSuite() {

}

func (a *AppConfigSuite) Test_GetConfig() {
	cfg1 := GetAppConfig()
	cfg2 := GetAppConfig()
	a.Assert().Equal(cfg1, cfg2)
}

func Test_MainConfig(t *testing.T) {
	suite.Run(t, &AppConfigSuite{})
}
