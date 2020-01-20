package dto

import (
	"math/big"
	"payments/src/types"
)

type PaymentReq struct {
	State 		types.TxState	`json:"state"`
	Amount 		*big.Float		`json:"amount"`
	TxId		string			`json:"transactionId"`
	SourceType 	types.TxSource	`json:"-"`
}
