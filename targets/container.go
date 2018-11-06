package targets

import (
	"flag"
	"fmt"
	"os"
)

func Startup() {
	fmt.Println(os.Args)

	playerName := flag.String("player", "", "eos account name")
	flag.Parse()

	if *playerName == "" {
		flag.Usage()
		os.Exit(0)
	}
	t := NewLLG(*playerName)
	go t.Run()
}
