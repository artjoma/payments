package data

import (
	"math/big"
	"payments/src/types"
	"time"
)

type TransactionData struct {
	Id 			int64
	TxId		string
	CreatedAt	time.Time
	PrevBalance	*big.Float
	Amount		*big.Float
	Balance		*big.Float
	State 		types.TxState
	Source		types.TxSource
	Status		types.TxStatus
}
