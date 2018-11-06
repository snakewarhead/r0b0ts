package main

import (
	"time"
	"github.com/snakewarhead/r0b0ts/utils"
	"github.com/snakewarhead/r0b0ts/targets"
	"github.com/snakewarhead/r0b0ts/services"
)

func main() {
	services.Startup()
	targets.Startup()

	utils.Logger.Info("startup ------------------------------------")

	// never stop
	for {
		time.Sleep(3 * time.Second)
	}
}