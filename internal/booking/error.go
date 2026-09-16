package booking

type ValidationError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e NotFoundError) Error() string { return e.Message }
