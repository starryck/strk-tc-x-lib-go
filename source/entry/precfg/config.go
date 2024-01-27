package precfg

import (
	"fmt"
	"path"
	"runtime"

	"github.com/caarlos0/env/v7"

	"github.com/forbot161602/x-lib-go/source/core/base/xbcfg"
	"github.com/forbot161602/x-lib-go/source/core/base/xbconst"
	"github.com/forbot161602/x-lib-go/source/core/toolkit/xbrand"
)

func NewConfig() *Config {
	config := (&builder{}).initialize().
		parseEnv().
		setBasePath().
		setServiceID().
		setServiceDeveloping().
		build()
	return config
}

var ServiceEnvironmentDevelopingMap = map[string]bool{
	xbconst.EnvironmentLocal: true,
	xbconst.EnvironmentDev:   true,
	xbconst.EnvironmentSIT:   true,
	xbconst.EnvironmentUAT:   false,
	xbconst.EnvironmentStage: false,
	xbconst.EnvironmentProd:  false,
}

func ParseEnv(config xbcfg.SpecConfig) {
	if err := env.Parse(config); err != nil {
		panic(err)
	}
}

func MakeBasePath(back int) string {
	_, dir, _, _ := runtime.Caller(1)
	for i := 0; i < back; i++ {
		dir = path.Dir(dir)
	}
	return dir
}

func MakeServiceID() string {
	return xbrand.MakeUUID4()
}

func MakeServiceDeveloping(srvEnv string) bool {
	if yes, ok := ServiceEnvironmentDevelopingMap[srvEnv]; ok {
		return yes
	} else {
		panic(fmt.Sprintf("Config does not support service environment `%s`.", srvEnv))
	}
}

type Config struct {
	BasePath  string `json:"basePath"`
	GitTag    string `json:"gitTag" env:"GIT_TAG,notEmpty"`
	GitCommit string `json:"gitCommit" env:"GIT_COMMIT,notEmpty"`

	ServiceID          string `json:"serviceID"`
	ServiceCode        string `json:"serviceCode" env:"SRV_CODE" envDefault:"S001"`
	ServiceName        string `json:"serviceName" env:"SRV_NAME" envDefault:"golang-lib"`
	ServicePort        int    `json:"servicePort" env:"SRV_PORT" envDefault:"80"`
	ServiceProject     string `json:"serviceProject" env:"SRV_PROJECT" envDefault:"open"`
	ServiceVersion     string `json:"serviceVersion" env:"SRV_VERSION" envDefault:"v1"`
	ServiceEnvironment string `json:"serviceEnvironment" env:"SRV_ENVIRONMENT,notEmpty"`
	ServiceLogLevel    string `json:"serviceLogLevel" env:"SRV_LOG_LEVEL" envDefault:"INFO"`
	ServiceTesting     bool   `json:"serviceTesting" env:"SRV_TESTING" envDefault:"false"`
	ServiceDebugging   bool   `json:"serviceDebugging" env:"SRV_DEBUGGING" envDefault:"false"`
	ServiceDeveloping  bool   `json:"serviceDeveloping"`
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

type builder struct {
	config *Config
}

func (builder *builder) build() *Config {
	return builder.config
}

func (builder *builder) initialize() *builder {
	builder.config = &Config{}
	return builder
}

func (builder *builder) parseEnv() *builder {
	ParseEnv(builder.config)
	return builder
}

func (builder *builder) setBasePath() *builder {
	builder.config.BasePath = MakeBasePath(4)
	return builder
}

func (builder *builder) setServiceID() *builder {
	builder.config.ServiceID = MakeServiceID()
	return builder
}

func (builder *builder) setServiceDeveloping() *builder {
	builder.config.ServiceDeveloping = MakeServiceDeveloping(builder.config.ServiceEnvironment)
	return builder
}
