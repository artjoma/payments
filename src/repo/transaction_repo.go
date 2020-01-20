package repo

import (
	"context"
	"math/big"
	"payments/src/data"
	"strconv"
	"strings"
)

type TransactionRepo struct{

}

func (repo *TransactionRepo) Save(tx *data.TransactionData) error{
	dbTx, err := connPool.Begin(context.Background())
	if err != nil{
		return err
	}
	amount, _ := tx.Amount.Float64()
	prevBalance, _ := tx.PrevBalance.Float64()
	balance, _ := tx.Balance.Float64()
	_, err = dbTx.Exec(context.Background(),
		"INSERT INTO transactions(tx_id, prev_balance, amount, balance, state, source, status) " +
		"	VALUES($1, $2, $3, $4, $5, $6, $7)",
		tx.TxId, prevBalance, amount, balance, tx.State, tx.Source, tx.Status)
	if err != nil{
		dbTx.Rollback(context.Background())
		return err
	}

	_, err = dbTx.Exec(context.Background(), "UPDATE accounts SET balance=$1, updated_at=now() WHERE id=1",
		balance)
	if err == nil{
		dbTx.Commit(context.Background())
	}else{
		dbTx.Rollback(context.Background())
		return err
	}

	return nil
}


func (repo *TransactionRepo) GetLastN (count int) ([]*data.TransactionData, error){
	rows, err := connPool.Query(context.Background(),
		"SELECT id, tx_id, created_at, prev_balance, amount, balance, state, source, status " +
		"	FROM transactions " +
		"	ORDER BY id " +
		"	DESC LIMIT $1", count)
	if err != nil{
		return nil, err
	}

	defer rows.Close()
	list := make([]*data.TransactionData, 0, count)
	for rows.Next(){
		entity := &data.TransactionData{}
		prevBalance := float64(0)
		amount := float64(0)
		balance := float64(0)
		if err = rows.Scan(&entity.Id, &entity.TxId, &entity.CreatedAt, &prevBalance, &amount, &balance,
			&entity.State, &entity.Source, &entity.Status); err != nil{
			return nil, err
		}
		entity.PrevBalance = big.NewFloat(prevBalance)
		entity.Amount = big.NewFloat(amount)
		entity.Balance = big.NewFloat(balance)
		list = append(list, entity)
	}

	return list, err
}

func (repo *TransactionRepo) UpdateTransactionStatus(list []*data.TransactionData, balance *big.Float) error{
	dbTx, err := connPool.Begin(context.Background())
	if err != nil{
		return err
	}
	ids := make([]string, len(list))
	for i, entity := range list {
		ids[i] = strconv.Itoa(int(entity.Id))
	}
	idsStr := strings.Join(ids[:], ",")

	_, err = dbTx.Exec(context.Background(), "UPDATE transactions SET status='canceled' WHERE id IN (" + idsStr +")")
	if err == nil{
		dbTx.Commit(context.Background())
	}else{
		dbTx.Rollback(context.Background())
		return err
	}
	balanceF, _ := balance.Float64()
	_, err = dbTx.Exec(context.Background(), "UPDATE accounts SET balance=$1, updated_at=now() WHERE id=1",
		balanceF)
	if err == nil{
		dbTx.Commit(context.Background())
	}else{
		dbTx.Rollback(context.Background())
		return err
	}

	return nil
}