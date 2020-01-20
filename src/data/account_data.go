package data

import (
	"math/big"
	"time"
)

type AccountData struct {
	Id 			int64
	Balance 	*big.Float
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}
