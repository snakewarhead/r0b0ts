package targets

import (
	"time"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	gameContract = "llgcontract1"
	gameTableURL = "https://api.lelego.io/gameTable"
)

type gameStatus int

const (
	wait gameStatus = iota
	start
	end
)

type gameTable struct {
	hand_id string `json:"hand_id"`

	curr_time int64 `json:"curr_time, string"`
	draw_time int64 `json:"draw_time, string"`

	player_bets      string `json:"player_bets"`
	player_pair_bets string `json:"player_pair_bets"`

	tie_bets string `json:"tie_bets"`

	banker_bets      string `json:"banker_bets"`
	banker_pair_bets string `json:"banker_pair_bets"`

	bets_id int64 `json:"bets_id"`
	frozen int64 `json:"frozen"`
	deck string `json:"deck"`
	bets string `json:"bets"`
}

type llg struct {
	Player      string
	gameStatus  gameStatus
	currentGame *gameTable
	gameHistory map[string]*gameTable
}

func NewLLG(player string) *llg {
	return &llg{
		Player:      player,
		gameStatus:  wait,
		currentGame: nil,
		gameHistory: make(map[string]*gameTable),
	}
}

func (t *llg) Run() {
	for {
		t.doRun()
		time.Sleep(60 * time.Microsecond)
	}
}

func (t *llg) doRun() {
	// need never sto5
	defer utils.RecoverAndLog("llg", "doRun")

	table := &gameTable{}
	if err := utils.HttpGet(gameTableURL, table); err != nil {
		utils.Logger.Error(err)
		return
	}
}
