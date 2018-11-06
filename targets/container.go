package targets

import (
	"flag"
	"fmt"
	"os"

	"github.com/snakewarhead/r0b0ts/targets/llg"
)

func Startup() {
	fmt.Println(os.Args)

	playerName := flag.String("player", "", "eos account name")
	flag.Parse()

	if *playerName == "" {
		flag.Usage()
		os.Exit(0)
	}
	t := llg.NewLLG(*playerName)
	go t.Run()
}
