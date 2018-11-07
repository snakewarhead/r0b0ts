package llg

import (
	"time"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	gameContract = "llgcontract1"
)

type playerState int

const (
	wait playerState = iota
	sitout
	betted
	win
	loss
)

type llg struct {
	Player      string
	gameHistory *gameHistory
}

func NewLLG(player string) *llg {
	return &llg{
		Player:      player,
		gameHistory: initGameHistory(),
	}
}

func (t *llg) Run() {
	for {
		t.doRun()
		time.Sleep(1000 * time.Millisecond)
	}
}

func (t *llg) doRun() {
	// need never stop
	defer utils.RecoverAndLog("llg", "doRun")

	if err := t.gameHistory.update(); err != nil {
		utils.Logger.Error(err)
		return
	}

	// fsm
	// switch t.gameHistory.peek().state {

	// case initTable:
	// 	_
	// }
}
