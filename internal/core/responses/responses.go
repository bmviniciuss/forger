package responses

type Response struct {
	StatusCode int           `json:"-"`
	Error      ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
}

func NewNotFoundResponse() *Response {
	return &Response{
		StatusCode: 404,
		Error: ErrorResponse{
			Code:    "not_found",
			Message: "Not found",
		},
	}
}

func NewInternalErrorResponse(message, reason string) *Response {
	return &Response{
		StatusCode: 500,
		Error: ErrorResponse{
			Code:    "internal_error",
			Message: message,
			Reason:  reason,
		},
	}
}
