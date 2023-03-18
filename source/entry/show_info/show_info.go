package show_info

import (
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gblog"
	"github.com/forbot161602/pbc-golang-lib/source/entry/config"
)

func Execute() error {
	gblog.WithFields(gblog.Fields{
		"gitTag":    config.GetGitTag(),
		"gitCommit": config.GetGitCommit(),
	}).Info("Log info message.")
	return nil
}
