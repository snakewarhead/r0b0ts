package llg

import (
	"time"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	eosContract = "eosio.token"
	gameAccount = "llgcontract1"
	historyMax   = 20
	defaultGameRestInterval = 1000 * time.Millisecond
)

type llg struct {
	player       *player
	gameHistory  *gameHistory
	restInterval time.Duration
}

func NewLLG(name string) *llg {
	return &llg{
		player:       newPlayer(name),
		gameHistory:  initGameHistory(),
		restInterval: defaultGameRestInterval,
	}
}

func (t *llg) Run() {
	for {
		t.doRunLucky()
		time.Sleep(t.restInterval)
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
	currentGame := t.gameHistory.peek()
	if currentGame == nil {
		return
	}

	switch currentGame.state {
	case initTable:
		// bet
	case drawingCards:
		// settle
	case endTable:
		// server is down, so do nothing
	}
}

func (t *llg) doRunLucky() {
	// need never stop
	defer utils.RecoverAndLog("llg", "doRun")

	if err := t.gameHistory.update(); err != nil {
		utils.Logger.Error(err)
		return
	}

	// fsm
	currentGame := t.gameHistory.peek()
	if currentGame == nil {
		return
	}

	switch currentGame.state {
	case initTable:
	case drawingCards:
		rest := currentGame.restTimeForBetting()
		utils.Logger.Debug("restTimeForBetting 1 -- %d", rest)

		if t.player.hasBetted(currentGame) {
			return
		}

		if rest < (int64)(2000000) {
			t.restInterval = 500 * time.Millisecond
		} else {
			// more than 2 sec, do nothing
			return
		}

		if rest < (int64)(1000000) {
			t.restInterval = 100 * time.Millisecond
		}

		if rest < (int64)(500000) {
			t.restInterval = 50 * time.Millisecond
		}
		utils.Logger.Debug("restTimeForBetting 2 -- %d", rest)

		pre := time.Now().UnixNano()
		currentGame.fetchGameResult()
		aft := time.Now().UnixNano()
		utils.Logger.Debug("fetchGameResult -- %d", (aft - pre) / (int64)(time.Millisecond))
		if currentGame.result == nil {
			return
		}
		// betting
		if err := t.player.betting(currentGame); err != nil {
			utils.Logger.Error(err)
		}
	case endTable:
		// server is down, so do nothing
		t.restInterval = defaultGameRestInterval
	}
}
