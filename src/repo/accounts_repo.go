package repo

import (
	"context"
	"math/big"
)

type AccountRepo struct {
}

func (repo *AccountRepo) GetBalance(accountId int64) (*big.Float, error){
	var balance float64 = 0.0
	connPool.QueryRow(context.Background(),
		"SELECT balance FROM accounts WHERE id=$1", accountId).Scan(&balance)
	return big.NewFloat(balance), nil
}
