package client

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"math/big"
	"math/rand"
	"net/http"
	"payments/src/api"
	"payments/src/dto"
	"payments/src/types"
	"sync"
	"testing"
	"time"
)
const(
	apiMethod = "http://127.0.0.1:8095/payment/new"
)
func init(){
	rand.Seed(time.Now().UnixNano())
}

func TestSimpleBot(t *testing.T){
	bots(t, 1, 100)
}

func TestHugeBot(t *testing.T){
	bots(t, 2, 1_000_000)
}

//very simple concurrent requests simulator
func bots(t *testing.T, threadCount, requestCount int){
	errorsCount := atomic.NewUint64(0)
	netTransport := &http.Transport{
		TLSHandshakeTimeout: 3 * time.Second,
	}
	httpClient := &http.Client{
		Timeout: time.Second * 8,
		Transport: netTransport,
	}

	clientPool := &sync.WaitGroup{}

	for i:= 0; i < threadCount; i++ {
		clientPool.Add(1)
		go func() {
			for j := 0; j < requestCount; j++ {
				requestModel := dto.PaymentReq{}
				requestModel.Amount = big.NewFloat(1 + rand.Float64() * (100.19 - 1))
				requestModel.TxId = uuid.New().String()
				if j % 2 == 0{
					requestModel.State = types.TxStateWin
				}else{
					requestModel.State = types.TxStateLost
				}

				data, err := json.Marshal(requestModel)

				request, _ := http.NewRequest("POST", apiMethod, bytes.NewBuffer(data))
				request.Header.Set("Content-type", "application/json")
				request.Header.Set(api.HeaderSourceType, string(types.TxSourceClient))
				response, err := httpClient.Do(request)
				if err == nil {
					if err != nil {
						t.Error("body err: " + err.Error())
						errorsCount.Inc()
						continue
					}
					respModel := dto.BaseResponse{}
					err = json.NewDecoder(response.Body).Decode(&respModel)
					if err != nil {
						t.Error("resp invalid err: " + err.Error())
						errorsCount.Inc()
						continue
					}
					if respModel.Err != nil {
						t.Error("resp err: " + respModel.Err.Message)
						errorsCount.Inc()
						continue
					}
				} else {
					t.Error("err: " + err.Error())
					errorsCount.Inc()
					continue
				}
			}
			clientPool.Done()
		}()
	}

	clientPool.Wait()
	log.Infof("errors count: %v", errorsCount)

	if errorsCount.Load() != uint64(0) {
		t.Fail()
	}
}