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
	handID string
	txID   string
	memo   string
}

type player struct {
	name       string
	betHistory map[string]*betInfo
}

func newPlayer(name string) *player {
	return &player{
		name:       name,
		betHistory: make(map[string]*betInfo),
	}
}

func (p *player) betting(currentGame *gameTable) error {
	memo := currentGame.HandID + ":player"
	txid, err := services.PushTransaction(eosContract, p.name, gameAccount, memo, "EOS", true, "1.0000", "")
	if err != nil {
		return err
	}
	utils.Logger.Info(txid)

	p.betHistory[currentGame.HandID] = &betInfo{currentGame.HandID, txid, memo}

	return err
}

func (p *player) hasBetted(currentGame *gameTable) bool {
	return p.betHistory[currentGame.HandID] != nil
}