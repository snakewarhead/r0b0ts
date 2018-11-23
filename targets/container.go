package targets

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/snakewarhead/r0b0ts/targets/llg"
)

func Startup() {
	fmt.Println(os.Args)

	playerName := flag.String("player", "", "eos account name")
	bettimes := flag.Int("bettimes", math.MaxInt32, "set the bet times, -1 is infinite")
	flag.Parse()

	if *playerName == "" {
		flag.Usage()
		os.Exit(0)
	}
	t := llg.NewLLG(*playerName, *bettimes)
	go t.Run()
}
