package structs

type ErrorInputNotSanitized struct{}

func (e *ErrorInputNotSanitized) Error() string {
	return "Input is not sanitized"
}
