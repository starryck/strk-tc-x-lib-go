package show_info

import (
	"github.com/forbot161602/x-lib-go/source/core/base/xbcfg"
	"github.com/forbot161602/x-lib-go/source/core/utility/xblogger"
)

func Execute() error {
	xblogger.WithFields(xblogger.Fields{
		"gitTag":    xbcfg.GetGitTag(),
		"gitCommit": xbcfg.GetGitCommit(),
	}).Info("Log info message.")
	return nil
}
