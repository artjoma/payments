package api

import (
	"math/big"
	"payments/src/dto"
	"payments/src/service"
	"payments/src/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	HeaderSourceType = "Source-Type"
)

var (
	validate *validator.Validate = validator.New()
)

func newPaymentHandler(ctx *gin.Context) {
	request := &dto.PaymentReq{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		invalidRequestErr("", "invalid json", ctx)
		return
	}
	request.SourceType = types.TxSource(ctx.GetHeader(HeaderSourceType))

	if err := validate.Struct(request); err != nil {
		invalidRequestErr(request.TxId, err.Error(), ctx)
		return
	}
	if !request.SourceType.IsValid() {
		invalidRequestErr(request.TxId, "invalid Source-Type header", ctx)
		return
	}
	if !request.State.IsValid() {
		invalidRequestErr(request.TxId, "invalid state field", ctx)
		return
	}
	if request.Amount.Cmp(big.NewFloat(0)) < 1 || request.Amount.IsInf() {
		invalidRequestErr(request.TxId, "invalid amount field", ctx)
		return
	}

	//request is valid, put to process queue
	service.GetPaymentService().AddNewPayment(request)
	successResponse("ACCEPTED", ctx)
}
