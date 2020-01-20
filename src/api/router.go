package api

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"payments/src/dto"
	"payments/src/types"
	"time"
)

var(
	httpServerInst *http.Server
)

//TODO Global PanicHandler
//func PanicHandler() {
//}

func ShutdownHttp() {
	log.Info("start Http server shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	httpServerInst.Shutdown(ctx)
	defer cancel()
	log.Info("end Http server shutdown")
}

func InitHttp(){
	//comment it for GIN access logging and etc.
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	httpServerInst = &http.Server{
		Addr:           ":" + os.Getenv("APP_HTTP_PORT"),
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 512,
		Handler:        buildRouter(),
		IdleTimeout:    10 * time.Minute,
	}

	go func() {
		if err := httpServerInst.ListenAndServe(); err != nil {
			log.Info("server shutdown ! " + err.Error())
		}
	}()
}

func invalidRequestErr(corrId, message string, ctx *gin.Context){
	appErr := &types.AppErr{}
	appErr.Id = types.ErrInvalidRequest
	appErr.Message = message
	appErr.CorrelationId = corrId
	resp := &dto.BaseResponse{}
	resp.Err = appErr
	ctx.JSON(http.StatusBadRequest, resp)
}

func successResponse(result interface{}, ctx *gin.Context){
	ctx.JSON(http.StatusOK, &dto.BaseResponse{Result: result})
}

func buildRouter() *gin.Engine{
	router := gin.Default()
	router.POST("payment/new", newPaymentHandler)

	return router
}

