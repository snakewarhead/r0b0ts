package llg

import (
	"github.com/snakewarhead/r0b0ts/services"
	"github.com/snakewarhead/r0b0ts/utils"
)

type playerState int

const (
	wait playerState = iota
	sitout
	betted
	win
	loss
)

type betInfo struct {
	handID  string
	txID    string
	memo    string
	success bool
}

type player struct {
	name       string
	betHistory *utils.LimitStackMap
}

func newPlayer(name string) *player {
	return &player{
		name:       name,
		betHistory: utils.NewLimitStackMap(historyMax),
	}
}

func (p *player) betting(currentGame *gameTable, memo, amount string) error {
	txid, err := services.PushTransaction(eosContract, p.name, gameAccount, memo, "EOS", true, amount, "")
	utils.Logger.Info("txid:%s, err:%v", txid, err)

	// betted yet if err occures
	p.betHistory.Push(currentGame.HandID, &betInfo{currentGame.HandID, txid, memo, err == nil})
	return err
}

func (p *player) hasBetted(currentGame *gameTable) bool {
	return p.betHistory.Find(currentGame.HandID) != nil
}
