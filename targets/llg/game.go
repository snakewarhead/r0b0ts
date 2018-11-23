package llg

import (
	"time"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	eosContract             = "eosio.token"
	gameAccount             = "llgcontract1"
	historyMax              = 20
	defaultGameRestInterval = 1000 * time.Millisecond
)

type llg struct {
	player       *player
	gameHistory  *gameHistory
	restInterval time.Duration
	currentGame  *gameTable
}

func NewLLG(name string) *llg {
	return &llg{
		player:       newPlayer(name),
		gameHistory:  initGameHistory(),
		restInterval: defaultGameRestInterval,
	}
}

func (t *llg) Run() {
	var now, last, delta int64
	now = time.Now().UnixNano()
	last = now
	for {
		now = time.Now().UnixNano()
		delta = now - last

		t.doRunLucky(delta)
		time.Sleep(t.restInterval)

		last = now
	}
}

func (t *llg) doRun(dt int64) {
	// need never stop
	defer utils.RecoverAndLog("llg", "doRun")

	if err := t.gameHistory.update(dt); err != nil {
		utils.Logger.Error(err)
		return
	}

	// fsm
	currentGame := t.gameHistory.peek()
	if currentGame == nil {
		return
	}
	// new game now
	if currentGame != t.currentGame {
		t.currentGame = currentGame
		t.restInterval = defaultGameRestInterval
	}
	utils.Logger.Debug("game state -- %d", t.currentGame.state)

	switch t.currentGame.state {
	case initTable:
		// bet
	case drawingCards:
		// settle
	case endTable:
		// server is down, so do nothing
	}
}

func (t *llg) doRunLucky(dt int64) {
	// need never stop
	defer utils.RecoverAndLog("llg", "doRun")

	if err := t.gameHistory.update(dt); err != nil {
		utils.Logger.Error(err)
		return
	}

	// fsm
	currentGame := t.gameHistory.peek()
	if currentGame == nil {
		return
	}
	// new game now
	if currentGame != t.currentGame {
		t.currentGame = currentGame
		t.restInterval = defaultGameRestInterval
	}
	// utils.Logger.Debug("game state -- %d", t.currentGame.state)

	switch t.currentGame.state {
	case initTable:
		// rest := t.currentGame.restTimeCorrectedForBetting()
		// utils.Logger.Debug("restTimeForBetting 1 -- %d, %d", rest, t.currentGame.restTimeForBetting())

		if t.player.hasBetted(t.currentGame) {
			return
		}

		rest := t.currentGame.restTimeCorrectedForBetting()
		if rest <= (int64)(3000000) {
			t.restInterval = 500 * time.Millisecond
		} else {
			// more than 2 sec, do nothing
			return
		}

		if rest <= (int64)(1000000) {
			t.restInterval = 100 * time.Millisecond
		}

		if rest <= (int64)(500000) {
			t.restInterval = 50 * time.Millisecond
		}
		utils.Logger.Debug("restTimeForBetting 2 -- %d", rest)

		pre := time.Now().UnixNano()
		t.currentGame.fetchGameResult()
		aft := time.Now().UnixNano()
		utils.Logger.Debug("fetchGameResult -- %d, %v", (aft-pre)/(int64)(time.Millisecond), t.currentGame.result)
		if t.currentGame.result == nil {
			return
		}
		// betting
		winners := t.currentGame.result.whoAreWinners()
		if len(winners) == 0 {
			return
		}
		memo := t.currentGame.HandID + ":" + winners[0]
		t.player.betting(t.currentGame, memo, "1.0000")

	case drawingCards:
		t.restInterval = defaultGameRestInterval
	case endTable:
		// server is down, so do nothing
		t.restInterval = defaultGameRestInterval
	}
}
