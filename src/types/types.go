package types

type TxState string
const(
	TxStateWin 	TxState = "win"
	TxStateLost TxState = "lost"
)
func (state TxState) IsValid() bool {
	switch state{
	case TxStateLost, TxStateWin:
		return true
	}
	return false
}

type TxSource string
const(
	TxSourceClient	TxSource = "client"
	TxSourceGame	TxSource = "game"
	TxSourceServer	TxSource = "server"
	TxSourcePayment	TxSource = "payment"
)

func (src TxSource) IsValid() bool {
	switch src {
	case TxSourceClient, TxSourceGame, TxSourceServer, TxSourcePayment:
		return true
	}
	return false
}


type TxStatus string
const(
	TxStatusCanceled 			TxStatus  = "canceled"
	TxStatusComplete			TxStatus  = "complete"
)

