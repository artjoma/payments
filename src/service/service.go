package service

var(
	paymentService 		*PaymentService
	accountCacheService *AccountsCacheService
)

func InitService() {
	accountCacheService = NewAccountCache()
	paymentService 		= NewPaymentService()
}

func ShutdownService() {
	paymentService.destroy()
}

func GetPaymentService() *PaymentService{
	return paymentService
}