package api

import (
	"github.com/gin-gonic/gin"
	"math/big"
	"payments/src/dto"
	"payments/src/service"
	"payments/src/types"
)

const(
	HeaderSourceType = "Source-Type"
)

func newPaymentHandler(ctx *gin.Context){
	request := &dto.PaymentReq{}
	if err := ctx.ShouldBindJSON(request); err != nil{
		invalidRequestErr("", "invalid json", ctx)
		return
	}
	if len(request.TxId) < 5{
		invalidRequestErr(request.TxId, "invalid transactionId field", ctx)
		return
	}
	request.SourceType = types.TxSource(ctx.GetHeader(HeaderSourceType))
	if !request.SourceType.IsValid(){
		invalidRequestErr(request.TxId, "invalid Source-Type header", ctx)
		return
	}
	if !request.State.IsValid(){
		invalidRequestErr(request.TxId, "invalid state field", ctx)
		return
	}

	if request.Amount.Cmp(big.NewFloat(0)) < 1 || request.Amount.IsInf(){
		invalidRequestErr(request.TxId, "invalid amount field", ctx)
		return
	}

	//request is valid, put to process queue
	service.GetPaymentService().AddNewPayment(request)
	successResponse("ACCEPTED", ctx)
}

