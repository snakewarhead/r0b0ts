package llg

import (
	"fmt"

	"github.com/snakewarhead/r0b0ts/utils"
)

const (
	llgURL   = "https://api.lelego.io/"
	tableURL = llgURL + "gameTable"
	gameURL  = llgURL + "game/"
)

type gameState int

const (
	unknown gameState = iota
	initTable
	drawingCards
	settle
)

type gameResult struct {
	ID          string        `json:"_id"`
	SeedHashed  string        `json:"seedHashed"`
	NowHandID   string        `json:"nowHandID"`
	NextHandID  string        `json:"nextHandID"`
	HandNumber  int           `json:"handNumber"`
	PlayerCards []interface{} `json:"playerCards"`

	PlayerWin   bool          `json:"playerWin"`
	BankerWin   bool          `json:"bankerWin"`
	TieWin  	bool          `json:"tieWin"`

	PlayerPair  bool          `json:"playerPair"`
	BankerPair  bool          `json:"bankerPair"`

	GameOutcome string        `json:"gameOutcome"`
	V           int           `json:"__v"`
}

func (r *gameResult) isInvalid() bool {
	empty := r.ID == "" || r.NowHandID == "" || r.NextHandID == ""
	nobodyWin := !(r.PlayerWin || r.BankerWin || r.TieWin)
	return empty || nobodyWin
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
}

func (t *gameTable) isInvalid() bool {
	return t.HandID == "" || t.CurrTime == 0 || t.DrawTime == 0
}

func (t *gameTable) updateState() {
	t.state = unknown
	if t.CurrTime < t.DrawTime {
		t.state = initTable
	} else {
		t.state = drawingCards

		// get result
		if t.result == nil {
			results := make([]*gameResult, 1)
			if err := utils.HttpGetVar(gameURL + t.HandID, &results); err != nil {
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
}

const historyMax = 20

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
		h.history[current.HandID] = table
		return
	}

	if h.size >= historyMax {
		ancient := h.gameIDs[0]
		h.gameIDs = h.gameIDs[1:]

		delete(h.history, ancient)
		h.size--
	}

	h.gameIDs = append(h.gameIDs, table.HandID)
	h.history[table.HandID] = table
	h.size++

	h.store(table)
}

func (h *gameHistory) update() error {
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
	current.updateState()

	return nil
}
