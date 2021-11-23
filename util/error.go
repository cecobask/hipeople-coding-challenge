package util

type RequestError struct {
	Status  int
	Message string
	Err     error
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}
