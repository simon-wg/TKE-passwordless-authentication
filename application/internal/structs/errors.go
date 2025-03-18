package structs

type ErrorInputNotSanitized struct {
	Message string
}

func (e *ErrorInputNotSanitized) Error() string {
	return e.Message
}
