package api

import "fmt"

var ErrConflict = NewErrResourceConflict("", "")

type ErrNotFound struct {
	Msg string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("a matching resource was not found: %s", e.Msg)
}

func NewErrNotFound(msg string) *ErrNotFound {
	return &ErrNotFound{Msg: msg}
}

type ErrInvalidRequest struct {
	Msg string
}

type ErrResourceConflict struct {
	Resource string
	ID       string
}

func (e ErrResourceConflict) Error() string {
	msg := "a matching resource already exist"
	if e.Resource != "" {
		msg = fmt.Sprintf("%s: %s", msg, e.Resource)
	}
	return msg
}

func (e ErrInvalidRequest) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Msg)
}

func NewInvalidErr(msg string) *ErrInvalidRequest {
	return &ErrInvalidRequest{Msg: msg}
}

func NewErrResourceConflict(resource, id string) ErrResourceConflict {
	return ErrResourceConflict{
		Resource: resource,
		ID:       id,
	}
}
