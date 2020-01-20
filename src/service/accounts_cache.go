package service

import (
	log "github.com/sirupsen/logrus"
	"math/big"
	"payments/src/repo"
)

//Account balances cache, currently store only one record
type AccountsCacheService struct {
	storage     map[int64]*big.Float //accountId: account balance
	accountRepo *repo.AccountRepo
}

func NewAccountCache() *AccountsCacheService {
	cache := &AccountsCacheService{}
	cache.accountRepo = repo.GetAccountRepo()
	//single record, preallocate size of map if multiple accounts!
	cacheStorage := make(map[int64]*big.Float)
	cache.storage = cacheStorage
	balance, err := cache.accountRepo.GetBalance(1)
	if err != nil{
		panic("account with id=1 not found! " + err.Error())
	}
	log.Infof("initial balance: %s", balance.String())
	cacheStorage[1] = balance
	return cache
}

func (cache *AccountsCacheService) GetAccountBalance(accountId int64) *big.Float{
	return cache.storage[accountId]
}

func (cache *AccountsCacheService) PutAccountBalance(accountId int64, balance *big.Float){
	cache.storage[accountId] = balance
}
