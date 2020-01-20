package types

import "fmt"

const (
	ErrServiceUnavailable	= iota
	ErrInvalidRequest
)

type AppErr struct {
	Id				uint16	`json:"id"`			// error id
	CorrelationId	string 	`json:"cId,omitempty"`		// request id, may be empty
	Message			string	`json:"msg"`		// short message
	Err				error 	`json:"-"`			// native error
}

func (err *AppErr) Error() string{
	return fmt.Sprintf("id:%d msg:%s cId:%s", err.Id, err.Message, err.CorrelationId)
}

func (err *AppErr) Unwrap() error{
	return err.Err
}