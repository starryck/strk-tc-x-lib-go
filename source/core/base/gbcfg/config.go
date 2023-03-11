package gbcfg

var mConfig Config

type Config interface {
	GetBasePath() string
	GetGitTag() string
	GetGitCommit() string
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
