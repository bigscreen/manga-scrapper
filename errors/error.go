package errors

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	keyOp      = "op"
	keyCode    = "code"
	keyMessage = "message"

	ErrorCodeNotFound   = "err_not_found"
	ErrorCodeValidation = "err_validation"
)

type Error struct {
	op      string
	code    string
	message string
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Error() string {
	var ta []string
	if len(e.op) > 0 {
		ta = append(ta, fmt.Sprintf("%s: %s", keyOp, e.op))
	}
	if len(e.code) > 0 {
		ta = append(ta, fmt.Sprintf("%s: %s", keyCode, e.code))
	}
	if len(e.message) > 0 {
		ta = append(ta, fmt.Sprintf("%s: %s", keyMessage, e.message))
	}
	return fmt.Sprintf("{%s}", strings.Join(ta, " | "))
}

func (e Error) HttpStatusCode() int {
	switch e.code {
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

type ErrorOpt struct {
	key   string
	value string
}

func WithOp(op string) ErrorOpt {
	return ErrorOpt{key: keyOp, value: op}
}

func WithCode(code string) ErrorOpt {
	return ErrorOpt{key: keyCode, value: code}
}

func WithMessage(msg string) ErrorOpt {
	return ErrorOpt{key: keyMessage, value: msg}
}

func WithNotFoundErrorCode() ErrorOpt {
	return ErrorOpt{key: keyCode, value: ErrorCodeNotFound}
}

func WithValidationErrorCode() ErrorOpt {
	return ErrorOpt{key: keyCode, value: ErrorCodeValidation}
}

func WithError(err error) ErrorOpt {
	if err != nil {
		return ErrorOpt{key: keyMessage, value: err.Error()}
	}
	return ErrorOpt{}
}

func New(opts ...ErrorOpt) error {
	e := Error{}
	for _, opt := range opts {
		if opt.key == keyOp {
			e.op = opt.value
		}
		if opt.key == keyCode {
			e.code = opt.value
		}
		if opt.key == keyMessage {
			e.message = opt.value
		}
	}
	return e
}

func FromError(err error) *Error {
	if e, ok := err.(Error); ok {
		return &e
	}
	return nil
}
