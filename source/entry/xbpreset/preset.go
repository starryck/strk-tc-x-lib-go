package xbpreset

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/forbot161602/x-lib-go/source/core/base/xbcfg"
	"github.com/forbot161602/x-lib-go/source/entry/xbprecfg"
)

func init() {
	if err := godotenv.Load(); err == nil {
		fmt.Println("[INFO] The .env file has been successfully loaded.")
	}
	xbcfg.SetConfig(xbprecfg.NewConfig())
}
