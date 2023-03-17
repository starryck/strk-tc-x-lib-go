package gbshow_info

import (
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gblog"
	"github.com/forbot161602/pbc-golang-lib/source/entry/gbconfig"
)

func Execute() error {
	gblog.WithFields(gblog.Fields{
		"gitTag":    gbconfig.GetGitTag(),
		"gitCommit": gbconfig.GetGitCommit(),
	}).Info("Log info message.")
	return nil
}
