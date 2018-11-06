package models

import (
	"database/sql"
	"fmt"

	"github.com/snakewarhead/r0b0ts/utils"
)

type TransactionStatus int
type TransactionDirection int

const (
	InitTransactionStatus      TransactionStatus = iota // 0
	ConfirmedTransactionStatus                          // 1
)

const (
	_                       TransactionDirection = iota
	InTransactionDirection                       // 1
	OutTransactionDirection                      // 2
)

// this model is a db model as well as a json object which is response in http
type Transaction struct {
	ID         int    `json:"id"`
	CoinID     int    `json:"coin_id"`
	Contract   string `json:"contract"`
	TXID       string `json:"txid"`
	IsMain     int    `json:"is_main"`
	Symbol     string `json:"symbol"`
	Direction  int    `json:"direction"`
	Status     int    `json:"status"`
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     string `json:"amount"`
	Fee        string `json:"fee"`
	Memo       string `json:"memo"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	BlockNum   int64  `json:"block_num"`
}

func SaveTransaction(coinID int, contract string, isMain bool, txID, symbol, from, to, memo string, amount, fee string, status TransactionStatus, direction TransactionDirection, blockNum int64) error {
	stmt, err := db.Prepare("INSERT INTO transaction_history(coin_id, contract, tx_id, is_main, symbol, from_address, to_address, amount, fee, memo, create_time, update_time, status, direction, block_num) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	timeStamp := utils.GetCurrentTimestamp()

	var isMainN int
	if isMain {
		isMainN = 1
	} else {
		isMainN = 0
	}

	_, err = stmt.Exec(coinID, contract, txID, isMainN, symbol, from, to, amount, fee, memo, timeStamp, timeStamp, int(status), int(direction), blockNum)
	if err != nil {
		return err
	}
	return nil
}

func SaveOrUpdateTransaction(trx *Transaction) error {
	trxDB, err := FindOneTransaction(trx.CoinID, trx.TXID)
	if err != nil {
		return err
	}
	if trxDB == nil {
		SaveTransaction(
			trx.CoinID,
			trx.Contract,
			trx.IsMain == 1,
			trx.TXID,
			trx.Symbol,
			trx.From,
			trx.To,
			trx.Memo,
			trx.Amount,
			trx.Fee,
			TransactionStatus(trx.Status),
			TransactionDirection(trx.Direction),
			trx.BlockNum,
		)
	} else {
		if trx.Contract != trxDB.Contract ||
			trx.Symbol != trxDB.Symbol ||
			trx.From != trxDB.From ||
			trx.To != trxDB.To ||
			trx.Amount != trxDB.Amount ||
			trx.Memo != trxDB.Memo {
			return fmt.Errorf("SaveOrUpdateTransaction trx is not the same. trx:%v, trxDB:%v", trx, trxDB)
		}

		stmt, err := db.Prepare("UPDATE transaction_history SET status=?, block_num=?, update_time=? WHERE coin_id = ? and tx_id = ?")
		if err != nil {
			return err
		}

		res, err := stmt.Exec(
			trx.Status,
			trx.BlockNum,
			utils.GetCurrentTimestamp(),
			trx.CoinID,
			trx.TXID,
		)
		if err != nil {
			return err
		}

		affect, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affect != 1 {
			return fmt.Errorf("UpdateTransaction affect num was not 1")
		}
	}

	return nil
}

func FindTransactions(coinid int, direction TransactionDirection, contract, symbol, account, memo string, pos, offset int) ([]*Transaction, error) {
	var (
		rows *sql.Rows
		err  error
	)

	// pos, offset is opposite of the cause of sql(limit offset)
	whereCause := " WHERE coin_id = ? and direction=? and contract=? and symbol=? and to_address=?"
	limitCause := " LIMIT ? OFFSET ?"
	if len(memo) == 0 {
		rows, err = db.Query(
			"SELECT * FROM transaction_history"+whereCause+" ORDER BY id DESC"+limitCause,
			coinid,
			int(direction),
			contract,
			symbol,
			account,
			offset,
			pos,
		)
	} else {
		whereCause += " and memo=?"
		rows, err = db.Query(
			"SELECT * FROM transaction_history"+whereCause+" ORDER BY id DESC"+limitCause,
			coinid,
			int(direction),
			contract,
			symbol,
			account,
			memo,
			offset,
			pos,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trxs := make([]*Transaction, 0)
	for rows.Next() {
		t := &Transaction{}

		err = rows.Scan(dbColumns(t)...)
		if err != nil {
			return nil, err
		}

		trxs = append(trxs, t)
	}
	return trxs, nil
}

// if found nothing, it would return nil transaction and nil error
func FindOneTransaction(coinid int, trxid string) (*Transaction, error) {
	trx := &Transaction{}

	row := db.QueryRow("SELECT * FROM transaction_history WHERE coin_id = ? and tx_id = ?", coinid, trxid)
	err := row.Scan(dbColumns(trx)...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return trx, nil
}
