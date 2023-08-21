package api

import "fmt"

type ErrInvalidRequest struct {
	Msg string
}

func (e ErrInvalidRequest) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Msg)
}

func NewInvalidErr(msg string) *ErrInvalidRequest {
	return &ErrInvalidRequest{Msg: msg}
}
