package xbcfg

var mConfig Config

type Config interface {
	GetBasePath() string
	GetGitTag() string
	GetGitCommit() string

	GetServiceID() string
	GetServiceCode() string
	GetServiceName() string
	GetServicePort() int
	GetServiceProject() string
	GetServiceVersion() string
	GetServiceEnvironment() string
	GetServiceLogLevel() string
	GetServiceTesting() bool
	GetServiceDebugging() bool
	GetServiceDeveloping() bool

	GetPostgresHost() string
	GetPostgresPort() string
	GetPostgresName() string
	GetPostgresUser() string
	GetPostgresPassword() string
}

func GetConfig() Config {
	if mConfig == nil {
		panic("Config hasn't been created.")
	}
	return mConfig
}

func SetConfig(config Config) {
	mConfig = config
}

// Base definition

func GetBasePath() string {
	return GetConfig().GetBasePath()
}

func GetGitTag() string {
	return GetConfig().GetGitTag()
}

func GetGitCommit() string {
	return GetConfig().GetGitCommit()
}

// Core definition

func GetServiceID() string {
	return GetConfig().GetServiceID()
}

func GetServiceCode() string {
	return GetConfig().GetServiceCode()
}

func GetServiceName() string {
	return GetConfig().GetServiceName()
}

func GetServicePort() int {
	return GetConfig().GetServicePort()
}

func GetServiceProject() string {
	return GetConfig().GetServiceProject()
}

func GetServiceVersion() string {
	return GetConfig().GetServiceVersion()
}

func GetServiceEnvironment() string {
	return GetConfig().GetServiceEnvironment()
}

func GetServiceLogLevel() string {
	return GetConfig().GetServiceLogLevel()
}

func GetServiceTesting() bool {
	return GetConfig().GetServiceTesting()
}

func GetServiceDebugging() bool {
	return GetConfig().GetServiceDebugging()
}

func GetServiceDeveloping() bool {
	return GetConfig().GetServiceDeveloping()
}

// Postgres server

func GetPostgresHost() string {
	return GetConfig().GetPostgresHost()
}

func GetPostgresPort() string {
	return GetConfig().GetPostgresPort()
}

func GetPostgresName() string {
	return GetConfig().GetPostgresName()
}

func GetPostgresUser() string {
	return GetConfig().GetPostgresUser()
}

func GetPostgresPassword() string {
	return GetConfig().GetPostgresPassword()
}
