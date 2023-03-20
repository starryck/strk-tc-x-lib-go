package preset

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joho/godotenv"

	"github.com/forbot161602/pbc-golang-lib/source/core/base/gbcfg"
	"github.com/forbot161602/pbc-golang-lib/source/entry/precfg"
)

func init() {
	if err := godotenv.Load(); err == nil {
		fmt.Println("[INFO] The .env file has been successfully loaded.")
	}
	rand.Seed(time.Now().UnixNano())
	gbcfg.SetConfig(precfg.NewConfig())
}
