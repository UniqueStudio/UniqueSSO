package pkg

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResp(err error) *Response {
	return &Response{
		Success: false,
		Message: err.Error(),
		Data:    "",
	}
}

func SuccessResp(data interface{}) *Response {
	return &Response{
		Success: true,
		Message: "ok",
		Data:    data,
	}
}
