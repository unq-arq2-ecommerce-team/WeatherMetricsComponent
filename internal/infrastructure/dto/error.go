package dto

type ErrorMessage struct {
	Message     string
	Description string
}

func NewErrorMessage(msg, desc string) *ErrorMessage {
	return &ErrorMessage{
		Message:     msg,
		Description: desc,
	}
}
