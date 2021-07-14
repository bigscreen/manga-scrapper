package contract

import "github.com/bigscreen/manga-scrapper/errors"

type Response struct {
	Success bool         `json:"success"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *ErrorData   `json:"error,omitempty"`
}

type ErrorData struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    &data,
	}
}

func NewErrorResponse(err error) Response {
	code := ""
	message := err.Error()
	if e := errors.FromError(err); e != nil {
		code = e.Code()
		message = e.Message()
	}

	return Response{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
		},
	}
}

func (r Response) String() string {
	return getJson(r)
}
