package repo

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

//shared variables
var(
	connPool 		*pgxpool.Pool
	transactionRepo *TransactionRepo
	accountRepo		*AccountRepo
)

func InitRepo(maxConCount, port, host, user, password, dbName string) {
	log.Info("start init repo")
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%s",
		user, password, host, port, dbName, maxConCount)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil{
		panic("invalid connection config parameters")
	}
	connPool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil{
		panic("init DB connection err:" + err.Error())
	}
	transactionRepo = &TransactionRepo{}

	log.Info("end init repo")
}

func ShutdownRepo() {
	log.Info("start shutdown pool")
	connPool.Close()
	log.Info("end shutdown pool")
}

func GetAccountRepo() *AccountRepo{
	return accountRepo
}

func GetTransactionRepo() *TransactionRepo{
	return transactionRepo
}


