package show_info

import (
	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbcfg"
	"github.com/forbot161602/pbc-golang-lib/source/core/utility/gblog"
)

func Execute() error {
	gblog.WithFields(gblog.Fields{
		"gitTag":    gbcfg.GetGitTag(),
		"gitCommit": gbcfg.GetGitCommit(),
	}).Info("Log info message.")
	return nil
}
