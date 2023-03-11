package gbpreset

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joho/godotenv"

	"github.com/forbot161602/pbc-golang-lib/source/entry/gbconfig"
)

func init() {
	if err := godotenv.Load(); err == nil {
		fmt.Println("[INFO] The .env file has been successfully loaded.")
	}
	gbconfig.SetConfig()
	rand.Seed(time.Now().UnixNano())
}
