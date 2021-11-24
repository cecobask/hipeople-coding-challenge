package util

type RequestError struct {
	Status  int
	Message string
	Err     error
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func NewRequestError(status int, message string, err error) *RequestError {
	return &RequestError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}
