package xbprecfg

import (
	"fmt"
	"path"
	"runtime"

	"github.com/caarlos0/env/v11"

	"github.com/starryck/strk-tc-x-lib-go/source/core/base/xbcfg"
	"github.com/starryck/strk-tc-x-lib-go/source/core/base/xbconst"
	"github.com/starryck/strk-tc-x-lib-go/source/core/toolkit/xbrand"
)

func NewConfig() *Config {
	config := (&configBuilder{}).initialize().
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

func ParseEnv(config xbcfg.Config) {
	if err := env.Parse(config); err != nil {
		panic(err)
	}
}

func MakeBasePath(back int) string {
	_, dir, _, _ := runtime.Caller(1)
	for range back {
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
	BasePath  string `json:"basePath" env:"-"`
	GitTag    string `json:"gitTag" env:"GIT_TAG"`
	GitCommit string `json:"gitCommit" env:"GIT_COMMIT"`

	ServiceID          string `json:"serviceID" env:"-"`
	ServiceCode        string `json:"serviceCode" env:"SRV_CODE" envDefault:"S001"`
	ServiceName        string `json:"serviceName" env:"SRV_NAME" envDefault:"lib-go"`
	ServicePort        int    `json:"servicePort" env:"SRV_PORT" envDefault:"80"`
	ServiceProject     string `json:"serviceProject" env:"SRV_PROJECT" envDefault:"x"`
	ServiceVersion     string `json:"serviceVersion" env:"SRV_VERSION" envDefault:"v1"`
	ServiceEnvironment string `json:"serviceEnvironment" env:"SRV_ENVIRONMENT" envDefault:"prod"`
	ServiceLogLevel    string `json:"serviceLogLevel" env:"SRV_LOG_LEVEL" envDefault:"info"`
	ServiceTesting     bool   `json:"serviceTesting" env:"SRV_TESTING" envDefault:"false"`
	ServiceDebugging   bool   `json:"serviceDebugging" env:"SRV_DEBUGGING" envDefault:"false"`
	ServiceDeveloping  bool   `json:"serviceDeveloping" env:"-"`

	PostgresHost     string `json:"postgresHost" env:"POSTGRES_HOST"`
	PostgresPort     string `json:"postgresPort" env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresName     string `json:"postgresName" env:"POSTGRES_NAME"`
	PostgresUser     string `json:"postgresUser" env:"POSTGRES_USER"`
	PostgresPassword string `json:"postgresPassword" env:"POSTGRES_PASSWORD"`
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

// Postgres server

func (config *Config) GetPostgresHost() string {
	return config.PostgresHost
}

func (config *Config) GetPostgresPort() string {
	return config.PostgresPort
}

func (config *Config) GetPostgresName() string {
	return config.PostgresName
}

func (config *Config) GetPostgresUser() string {
	return config.PostgresUser
}

func (config *Config) GetPostgresPassword() string {
	return config.PostgresPassword
}

type configBuilder struct {
	config *Config
}

func (builder *configBuilder) build() *Config {
	return builder.config
}

func (builder *configBuilder) initialize() *configBuilder {
	builder.config = &Config{}
	return builder
}

func (builder *configBuilder) parseEnv() *configBuilder {
	ParseEnv(builder.config)
	return builder
}

func (builder *configBuilder) setBasePath() *configBuilder {
	builder.config.BasePath = MakeBasePath(4)
	return builder
}

func (builder *configBuilder) setServiceID() *configBuilder {
	builder.config.ServiceID = MakeServiceID()
	return builder
}

func (builder *configBuilder) setServiceDeveloping() *configBuilder {
	builder.config.ServiceDeveloping = MakeServiceDeveloping(builder.config.ServiceEnvironment)
	return builder
}
