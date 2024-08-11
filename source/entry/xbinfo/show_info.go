package xbinfo

import (
	"github.com/starryck/x-lib-go/source/core/base/xbcfg"
	"github.com/starryck/x-lib-go/source/core/utility/xblogger"
)

func Execute() error {
	xblogger.WithFields(xblogger.Fields{
		"config": xbcfg.GetConfig(),
	}).Info("Log info message.")
	return nil
}
