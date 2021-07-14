package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

type ErrorTestSuite struct {
	suite.Suite
}

func (s *ErrorTestSuite) TestGetCode() {
	c := "code"
	s.Equal(c, Error{code: c}.Code())
}

func (s *ErrorTestSuite) TestGetMessage() {
	m := "message"
	s.Equal(m, Error{message: m}.Message())
}

func (s *ErrorTestSuite) TestGetError() {
	cases := []struct {
		name     string
		op       string
		code     string
		message  string
		expected string
	}{
		{
			name:     "WhenOnlyOpIsSet",
			op:       "op",
			expected: "{op: op}",
		},
		{
			name:     "WhenOnlyCodeIsSet",
			code:     "code",
			expected: "{code: code}",
		},
		{
			name:     "WhenOnlyMessageIsSet",
			message:  "message",
			expected: "{message: message}",
		},
		{
			name:     "WhenAllFieldsAreSet",
			op:       "op",
			code:     "code",
			message:  "message",
			expected: "{op: op | code: code | message: message}",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			err := Error{op: tc.op, code: tc.code, message: tc.message}
			s.Equal(tc.expected, err.Error())
		})
	}
}

func (s *ErrorTestSuite) TestGetHttpStatusCode() {
	cases := []struct {
		name     string
		code     string
		expected int
	}{
		{
			name:     "WhenCodeIsErrorNotFound",
			code:     ErrorCodeNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "WhenCodeIsErrorValidation",
			code:     ErrorCodeValidation,
			expected: http.StatusBadRequest,
		},
		{
			name:     "WhenCodeIsEmpty",
			expected: http.StatusInternalServerError,
		},
		{
			name:     "WhenCodeIsNotListed",
			code:     "code",
			expected: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			err := Error{code: tc.code}
			s.Equal(tc.expected, err.HttpStatusCode())
		})
	}
}

func (s *ErrorTestSuite) TestNewError() {
	cases := []struct {
		name     string
		opts     []ErrorOpt
		expected Error
	}{
		{
			name:     "WhenNoOptIsSet",
			expected: Error{},
		},
		{
			name:     "WhenOpOptIsSet",
			opts:     []ErrorOpt{WithOp("op")},
			expected: Error{op: "op"},
		},
		{
			name:     "WhenCodeOptIsSet",
			opts:     []ErrorOpt{WithCode("code")},
			expected: Error{code: "code"},
		},
		{
			name:     "WhenMessageOptIsSet",
			opts:     []ErrorOpt{WithMessage("message")},
			expected: Error{message: "message"},
		},
		{
			name:     "WhenNotFoundErrorCodeOptIsSet",
			opts:     []ErrorOpt{WithNotFoundErrorCode()},
			expected: Error{code: ErrorCodeNotFound},
		},
		{
			name:     "WhenValidationErrorCodeOptIsSet",
			opts:     []ErrorOpt{WithValidationErrorCode()},
			expected: Error{code: ErrorCodeValidation},
		},
		{
			name:     "WhenErrorOptIsSet",
			opts:     []ErrorOpt{WithError(errors.New("some error"))},
			expected: Error{message: "some error"},
		},
		{
			name:     "WhenErrorOptWithNilIsSet",
			opts:     []ErrorOpt{WithError(nil)},
			expected: Error{},
		},
		{
			name:     "WhenMultipleOptAreSet",
			opts:     []ErrorOpt{WithOp("op"), WithCode("code"), WithMessage("message")},
			expected: Error{op: "op", code: "code", message: "message"},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			err, ok := New(tc.opts...).(Error)
			s.True(ok)
			s.Equal(tc.expected, err)
		})
	}
}

func (s *ErrorTestSuite) TestFromError() {
	cases := []struct {
		name     string
		err      error
		expected *Error
	}{
		{
			name:     "WhenCastedFromCustomError",
			err:      New(WithValidationErrorCode()),
			expected: &Error{code: ErrorCodeValidation},
		},
		{
			name:     "WhenNotCastedFromCustomError",
			err:      errors.New("foo"),
			expected: nil,
		},
		{
			name:     "WhenNilIsSet",
			expected: nil,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.Equal(tc.expected, FromError(tc.err))
		})
	}
}
