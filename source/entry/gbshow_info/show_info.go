package gbshow_info

import (
	"fmt"

	"github.com/forbot161602/pbc-golang-lib/source/entry/gbconfig"
)

func Execute() error {
	fmt.Println(map[string]string{
		"gitTag":    gbconfig.GetGitTag(),
		"gitCommit": gbconfig.GetGitCommit(),
	})
	return nil
}
