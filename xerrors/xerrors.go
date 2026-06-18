package xerrors

import (
	"errors"
	"fmt"
)

type Code int

const (
	Internal Code = iota + 1
	HIDInit
	DeviceNotFound
	DeviceOpen
	DeviceWrite
	ConfigHomeDir
	ConfigSave
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

func New(code Code, msg string) *Error { return &Error{Code: code, Message: msg} }

func Wrap(code Code, msg string, err error) *Error { return &Error{Code: code, Message: msg, Err: err} }

func CodeOf(err error) Code {
	if xe, ok := errors.AsType[*Error](err); ok {
		return xe.Code
	}
	return 0
}
