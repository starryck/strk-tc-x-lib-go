package precfg

import (
	"fmt"
	"path"
	"runtime"

	"github.com/caarlos0/env/v7"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbconst"
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbidtf"
)

var mConfig *Config

type Config struct {
	BasePath  string
	GitTag    string `env:"GIT_TAG,notEmpty"`
	GitCommit string `env:"GIT_COMMIT,notEmpty"`

	ServiceID          string
	ServiceCode        string `env:"SRV_CODE" envDefault:"S001"`
	ServiceName        string `env:"SRV_NAME" envDefault:"golang-lib"`
	ServicePort        int    `env:"SRV_PORT" envDefault:"80"`
	ServiceProject     string `env:"SRV_PROJECT" envDefault:"open"`
	ServiceVersion     string `env:"SRV_VERSION" envDefault:"v1"`
	ServiceEnvironment string `env:"SRV_ENVIRONMENT,notEmpty"`
	ServiceLogLevel    string `env:"SRV_LOG_LEVEL" envDefault:"INFO"`
	ServiceTesting     bool   `env:"SRV_TESTING" envDefault:"false"`
	ServiceDebugging   bool   `env:"SRV_DEBUGGING" envDefault:"false"`
	ServiceDeveloping  bool
}

var EnvironmentDevelopingMap = map[string]bool{
	gbconst.EnvironmentLocal: true,
	gbconst.EnvironmentDev:   true,
	gbconst.EnvironmentSIT:   true,
	gbconst.EnvironmentUAT:   false,
	gbconst.EnvironmentStage: false,
	gbconst.EnvironmentProd:  false,
}

func (config *Config) Build() *Config {
	config.ParseEnv().
		SetBasePath().
		SetServiceID().
		SetServiceDeveloping()
	return config
}

func (config *Config) ParseEnv() *Config {
	if err := env.Parse(config); err != nil {
		panic(err)
	}
	return config
}

func (config *Config) SetBasePath() *Config {
	_, file, _, _ := runtime.Caller(0)
	config.BasePath = path.Dir(path.Dir(path.Dir(path.Dir(file))))
	return config
}

func (config *Config) SetServiceID() *Config {
	config.ServiceID = gbidtf.MakeUUID4()
	return config
}

func (config *Config) SetServiceDeveloping() *Config {
	srvEnv := config.ServiceEnvironment
	if yes, ok := EnvironmentDevelopingMap[srvEnv]; ok {
		config.ServiceDeveloping = yes
	} else {
		panic(fmt.Sprintf("Config does not support service environment `%s`.", srvEnv))
	}
	return config
}

// Base definition

func (config *Config) GetBasePath() string {
	return config.BasePath
}

func (config *Config) GetGitTag() string {
	return config.GitTag
}

func (config *Config) GetGitCommit() string {
	return config.GitCommit
}

// Core definition

func (config *Config) GetServiceID() string {
	return config.ServiceID
}

func (config *Config) GetServiceCode() string {
	return config.ServiceCode
}

func (config *Config) GetServiceName() string {
	return config.ServiceName
}

func (config *Config) GetServicePort() int {
	return config.ServicePort
}

func (config *Config) GetServiceProject() string {
	return config.ServiceProject
}

func (config *Config) GetServiceVersion() string {
	return config.ServiceVersion
}

func (config *Config) GetServiceEnvironment() string {
	return config.ServiceEnvironment
}

func (config *Config) GetServiceLogLevel() string {
	return config.ServiceLogLevel
}

func (config *Config) GetServiceTesting() bool {
	return config.ServiceTesting
}

func (config *Config) GetServiceDebugging() bool {
	return config.ServiceDebugging
}

func (config *Config) GetServiceDeveloping() bool {
	return config.ServiceDeveloping
}
