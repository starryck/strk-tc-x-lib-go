package gbconfig

import (
	"path"
	"runtime"

	"github.com/caarlos0/env/v7"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbcfg"
)

var mConfig *Config

type Config struct {
	BasePath  string
	GitTag    string `env:"GIT_TAG,required"`
	GitCommit string `env:"GIT_COMMIT,required"`
}

func GetConfig() *Config {
	if mConfig == nil {
		panic("Config hasn't been created.")
	}
	return mConfig
}

func SetConfig() {
	if mConfig == nil {
		mConfig = newConfig()
	}
	return
}

func newConfig() *Config {
	config := (&Config{}).
		initialize().
		setConfig().
		setBasePath()
	return config
}

func (config *Config) initialize() *Config {
	if err := env.Parse(config); err != nil {
		panic(err)
	}
	return config
}

func (config *Config) setConfig() *Config {
	gbcfg.SetConfig(config)
	return config
}

func (config *Config) setBasePath() *Config {
	_, modulePath, _, _ := runtime.Caller(0)
	config.BasePath = path.Dir(path.Dir(path.Dir(modulePath)))
	return config
}

// Base definition

func GetBasePath() string {
	return GetConfig().GetBasePath()
}

func (config *Config) GetBasePath() string {
	return config.BasePath
}

func GetGitTag() string {
	return GetConfig().GetGitTag()
}

func (config *Config) GetGitTag() string {
	return config.GitTag
}

func GetGitCommit() string {
	return GetConfig().GetGitCommit()
}

func (config *Config) GetGitCommit() string {
	return config.GitCommit
}
