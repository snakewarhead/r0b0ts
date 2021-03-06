package llg

import (
	"fmt"
	"time"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	llgURL   = "https://api.lelego.io/"
	tableURL = llgURL + "gameTable"
	gameURL  = llgURL + "game/"

	gameDurationMax = 70 * time.Second
)

type gameState int

const (
	unknown gameState = iota
	initTable
	drawingCards
	endTable
)

type gameResult struct {
	ID          string        `json:"_id"`
	SeedHashed  string        `json:"seedHashed"`
	NowHandID   string        `json:"nowHandID"`
	NextHandID  string        `json:"nextHandID"`
	HandNumber  int           `json:"handNumber"`
	PlayerCards []interface{} `json:"playerCards"`

	PlayerWin bool `json:"playerWin"`
	BankerWin bool `json:"bankerWin"`
	TieWin    bool `json:"tieWin"`

	PlayerPair bool `json:"playerPair"`
	BankerPair bool `json:"bankerPair"`

	GameOutcome string `json:"gameOutcome"`
	V           int    `json:"__v"`
}

func (r *gameResult) isInvalid() bool {
	empty := r.ID == "" || r.NowHandID == "" || r.NextHandID == ""
	nobodyWin := !(r.PlayerWin || r.BankerWin || r.TieWin)
	return empty || nobodyWin
}

func (r *gameResult) whoAreWinners() []string {
	winners := make([]string, 0, 3)
	if r.PlayerWin {
		winners = append(winners, "player")
	}
	if r.BankerWin {
		winners = append(winners, "banker")
	}
	if r.TieWin {
		winners = append(winners, "tie")
	}
	if r.PlayerPair {
		winners = append(winners, "playerpair")
	}
	if r.BankerPair {
		winners = append(winners, "bankerpair")
	}
	return winners
}

type gameTable struct {
	HandID   string `json:"hand_id"`
	CurrTime int64  `json:"curr_time,string"`
	DrawTime int64  `json:"draw_time,string"`

	PlayerBets     string `json:"player_bets"`
	PlayerPairBets string `json:"player_pair_bets"`

	TieBets string `json:"tie_bets"`

	BankerBets     string `json:"banker_bets"`
	BankerPairBets string `json:"banker_pair_bets"`

	BetsID int64 `json:"bets_id"`
	Frozen int64 `json:"frozen"`

	Deck string        `json:"deck"`
	Bets []interface{} `json:"bets"`

	state  gameState   `json:"-"`
	result *gameResult `json:"-"`

	startTimeInServer int64 `json:"-"`
	elapseInServer    int64 `json:"-"`

	restTimeCorrected int64 `json:"-"`
}

func (t *gameTable) isInvalid() bool {
	return t.HandID == "" || t.CurrTime == 0 || t.DrawTime == 0
}

// @return micro second until drawing card
func (t *gameTable) restTimeForBetting() int64 {
	return t.DrawTime - t.CurrTime
}

// This is a little bit more precise
func (t *gameTable) restTimeCorrectedForBetting() int64 {
	return t.restTimeCorrected - t.elapseInServer/1000
}

func (t *gameTable) updateState(dt int64) {
	t.state = unknown
	if t.CurrTime < t.DrawTime {
		t.state = initTable
	} else {
		t.state = drawingCards
		t.fetchGameResult()
	}

	// http will consume the dt
	// t.restTimeCorrected -= dt / 1000
	// utils.Logger.Debug(t.restTimeCorrected, t.restTimeForBetting())
}

func (t *gameTable) fetchGameResult() {
	if t.result == nil {
		results := make([]*gameResult, 1)
		if err := utils.HttpGetVar(gameURL+t.HandID, &results); err != nil {
			utils.Logger.Error(err)
			return
		}
		if len(results) == 0 || results[0].isInvalid() {
			utils.Logger.Error("game result is invalid")
			return
		}
		t.result = results[0]
	}
}

type gameHistory struct {
	size int

	gameIDs []string
	history map[string]*gameTable
}

func initGameHistory() *gameHistory {
	return &gameHistory{
		size:    0,
		gameIDs: make([]string, 0, historyMax),
		history: make(map[string]*gameTable),
	}
}

func (h *gameHistory) peek() *gameTable {
	if h.size > 0 {
		return h.history[h.gameIDs[h.size-1]]
	}
	return nil
}

func (h *gameHistory) store(table *gameTable) {
	// TODO: save into db
}

func (h *gameHistory) push(table *gameTable) {
	// is newer
	current := h.peek()
	if current != nil && current.HandID == table.HandID {
		// update content
		h.history[current.HandID].CurrTime = table.CurrTime
		h.history[current.HandID].DrawTime = table.DrawTime

		h.history[current.HandID].PlayerBets = table.PlayerBets
		h.history[current.HandID].PlayerPairBets = table.PlayerPairBets
		h.history[current.HandID].TieBets = table.TieBets
		h.history[current.HandID].BankerBets = table.BankerBets
		h.history[current.HandID].BankerPairBets = table.BankerPairBets

		h.history[current.HandID].Bets = table.Bets
		h.history[current.HandID].Frozen = table.Frozen
		return
	}

	if h.size >= historyMax {
		ancient := h.gameIDs[0]
		h.gameIDs = h.gameIDs[1:]

		delete(h.history, ancient)
		h.size--
	}

	// calculate the correction of the rest time, it will be updated by delta time
	table.restTimeCorrected = table.restTimeForBetting()

	h.gameIDs = append(h.gameIDs, table.HandID)
	h.history[table.HandID] = table
	h.size++

	// start time
	table.startTimeInServer = time.Now().UnixNano()

	h.store(table)
}

func (h *gameHistory) updateHistoryState() {
	now := time.Now().UnixNano()
	for i := h.size - 1; i >= 0; i-- {
		t := h.history[h.gameIDs[i]]
		if t.state == endTable {
			continue
		}

		t.elapseInServer = now - t.startTimeInServer
		if t.elapseInServer > (int64)(gameDurationMax) {
			t.state = endTable
		}
	}
}

func (h *gameHistory) update(dt int64) error {
	// update history state, history must be history, avoid failures to fetch current then stuck the program
	h.updateHistoryState()

	// fetch current game table
	table := &gameTable{}
	if err := utils.HttpGetVar(tableURL, table); err != nil {
		return err
	}
	if table == nil || table.isInvalid() {
		return fmt.Errorf("game table info is invalid")
	}
	h.push(table)

	// update game table state
	current := h.peek()
	if current == nil {
		return fmt.Errorf("this is inpossible")
	}
	current.updateState(dt)

	return nil
}
