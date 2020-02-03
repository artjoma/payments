package dto

import (
	"math/big"
	"payments/src/types"
)

type PaymentReq struct {
	State      types.TxState  `json:"state", validate:"required"`
	Amount     *big.Float     `json:"amount", validate:"required"`
	TxId       string         `json:"transactionId", validate:"required"`
	SourceType types.TxSource `json:"-"`
}
