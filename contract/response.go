package contract

type Response struct {
	Success      bool         `json:"success"`
	Data         *interface{} `json:"data,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    &data,
	}
}

func NewErrorResponse(err error) Response {
	return Response{
		Success:      false,
		ErrorMessage: err.Error(),
	}
}

func (r Response) String() string {
	return getJson(r)
}
