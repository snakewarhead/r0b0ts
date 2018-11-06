package services

import (
	"time"
	"github.com/snakewarhead/r0b0ts/models"
	"github.com/snakewarhead/r0b0ts/utils"
)

var (
	nodeManager = make(anodeManager)
	nodeCurrent inode
	coin        *models.Coin
)

type inode interface {
	id() int
	bind(c *models.Coin)
	getBind() *models.Coin
	pushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error)
	getBalance(contract, account, symbol string) (string, error)
	obversing()
}

type anodeManager map[int]inode

func init() {
	// register all nodes
	n := &nodeEOS{}
	nodeManager[n.id()] = n

	// TODO: add others
}

func Startup() {
	var err error
	coin, err = models.GetCoinEnabled()
	if err != nil {
		utils.Logger.Critical("must have one enabled coin! %v", err)
		panic(err)
	}

	// find the node
	nodeCurrent = nodeManager[coin.ID]
	nodeCurrent.bind(coin)

	// goroutine for observing the transfer in block
	// go obverseTransactionsInChain()
}

func PushTransaction(contract, from, to, memo, symbol string, isMain bool, amount, fee string) (string, error) {
	utils.Logger.Info("PushTransaction 1 ----------------- contract:%s, from:%s, to:%s, memo:%s, symbol:%s, isMain:%t, amount:%s, fee:%s",
		contract,
		from,
		to,
		memo,
		symbol,
		isMain,
		amount,
		fee)

	txid, err := nodeCurrent.pushTransaction(contract, from, to, memo, symbol, isMain, amount, fee)
	utils.Logger.Info("PushTransaction 2 ----------------- txid:%s, err:%v", txid, err)
	if err != nil {
		return "", err
	}

	// persistent
	// note that this operation maybe action behind by scan goroutine, but don't care about because txid would be returned in any case
	errPersistent := models.SaveTransaction(
		nodeCurrent.getBind().ID,
		contract,
		isMain,
		txid,
		symbol,
		from,
		to,
		memo,
		amount,
		fee,
		models.InitTransactionStatus,
		models.OutTransactionDirection,
		0,
	)
	if errPersistent != nil {
		utils.Logger.Error("PushTransaction persistent ----------------- txid:%s, err:%v", txid, errPersistent)
	}

	// transaction is success, it must be response, ignore the other errors
	return txid, nil
}

func GetBalance(contract, account, symbol string) (string, error) {
	return nodeCurrent.getBalance(contract, account, symbol)
}

func GetTransactionsFromDB(direction models.TransactionDirection, contract, symbol, account, memo string, pos, offset int) ([]*models.Transaction, error) {
	return models.FindTransactions(coin.ID, direction, contract, symbol, account, memo, pos, offset)
}

func GetOneTransactionFromDB(trxid string) (*models.Transaction, error) {
	return models.FindOneTransaction(coin.ID, trxid)
}

func obverseTransactionsInChain() {
	for {
		// obverse block and find transactions about ours
		nodeCurrent.obversing()

		// sleep a while
		time.Sleep(30 * time.Second)
	}
}
