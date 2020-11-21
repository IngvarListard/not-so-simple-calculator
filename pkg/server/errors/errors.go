package errors

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}
