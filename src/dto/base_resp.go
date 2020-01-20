package dto

import "payments/src/types"

type BaseResponse struct{
	Result 	interface{}		`json:"result"`
	Err 	*types.AppErr	`json:"err"`
}

