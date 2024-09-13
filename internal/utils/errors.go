package utils

type errorMessage struct {
	Error string
}

func HandleError(err error) errorMessage {
	return errorMessage{
		Error: err.Error(),
	}
}
