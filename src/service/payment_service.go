package service

import (
	"github.com/jackc/pgconn"
	log "github.com/sirupsen/logrus"
	"math/big"
	"payments/src/data"
	"payments/src/dto"
	"payments/src/repo"
	"payments/src/types"
	"time"
)
/**
	Payment service save & cancel operations
 */
const(
	lastTransactionsCount = 10
	paymentBufferSize 	  = 256
	cancelTicker		  = time.Minute
)
type PaymentService struct {
	active          bool
	requestQueue    chan 	*dto.PaymentReq
	cancelTicker    *time.Ticker
	transactionRepo *repo.TransactionRepo
}

func NewPaymentService() *PaymentService{
	service := &PaymentService{}
	service.active = true
	service.requestQueue = make(chan *dto.PaymentReq, paymentBufferSize)
	service.cancelTicker = time.NewTicker(cancelTicker) //every minute some transactions will be canceled
	service.transactionRepo = repo.GetTransactionRepo()
	go service.requestQueueListener()
	return service
}

func (srv *PaymentService) destroy(){
	for len(srv.requestQueue) > 0 {
		time.Sleep(time.Millisecond * 100)
	}
	srv.active = false
	close(srv.requestQueue)
	srv.cancelTicker.Stop()
}

func (srv *PaymentService) AddNewPayment(payment *dto.PaymentReq){
	log.Infof("new payment: %s %s %s %f", payment.TxId, payment.SourceType, payment.State,
		payment.Amount)
	srv.requestQueue <- payment
}

func (srv *PaymentService) requestQueueListener(){
	for srv.active{
		select {
			//process payment
			case payment, ok := <-srv.requestQueue:
				if ok {
					srv.processPayment(payment)
				}
			//cancel tx timer
			case _, ok := <- srv.cancelTicker.C:
				if ok{
					srv.cancelTransactions()
				}
		}
	}
}

func (srv *PaymentService) cancelTransactions() {
	log.Info("start cancel txs")
	txList, err := srv.transactionRepo.GetLastN(lastTransactionsCount)
	if err != nil{
		log.Errorf("err: %s", err.Error())
		return
	}
	log.Infof("tx count: %d", len(txList))
	if len(txList) > 0 {
		balance := accountCacheService.GetAccountBalance(1)
		log.Infof("balance: %s", balance.String())
		var newBalance *big.Float = nil

		for _, tx := range txList {
			log.Infof("txId: %s id: %d %s %f", tx.TxId, tx.Id, tx.State, tx.Amount)
			if tx.Status == types.TxStatusComplete {
				tx.Status = types.TxStatusCanceled
				newBalance = tx.PrevBalance
			} else {
				log.Infof("ignored, status %s", tx.Status)
			}
		}
		if newBalance != nil {
			log.Infof("new balance: %s", newBalance.String())
			srv.transactionRepo.UpdateTransactionStatus(txList, newBalance)
			accountCacheService.PutAccountBalance(1, newBalance)
		}
	}

	log.Info("end cancel txs")
}

func (srv *PaymentService) processPayment(payment *dto.PaymentReq){
	start := time.Now()
	balance := accountCacheService.GetAccountBalance(1)
	amount := payment.Amount

	if payment.State == types.TxStateLost {
		if balance.Cmp(amount) < 0 {
			//TODO write rejected tx to somewhere
			log.Errorf("overdraft not supported, transaction rejected! balance: %s", balance.String())
			return
		}
		amount.Neg(amount)
	}

	newBalance := big.NewFloat(0)
	newBalance.Add(balance, amount)

	txData := &data.TransactionData{}
	txData.TxId = payment.TxId
	txData.Amount = amount
	txData.State = payment.State
	txData.Source = payment.SourceType
	txData.Status = types.TxStatusComplete
	txData.PrevBalance = balance
	txData.Balance = newBalance

	if err := srv.writeWithAttempt(txData); err != nil{
		return
	}

	accountCacheService.PutAccountBalance(1, newBalance)

	log.Infof("[%s] new balance: %s, q size: %d, took: %dµs", txData.TxId, newBalance.String(),
		len(srv.requestQueue), time.Since(start).Microseconds())
}

func (srv *PaymentService) writeWithAttempt(txData *data.TransactionData) error{
	if err := srv.transactionRepo.Save(txData); err != nil {
		pgErr := err.(*pgconn.PgError)
		//duplicate tx with same tx_id
		if pgErr.Code == "23505" {
			// TODO write to somewhere this tx
			log.Errorf("duplicate tx with id: %s", txData.TxId)
			return err
		}else{
			log.Error("err: %s", err.Error() + " . Try again...")
			time.Sleep(time.Second  * 2)
			//try again
			srv.writeWithAttempt(txData)
		}
	}

	return nil
}

