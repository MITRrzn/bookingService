package booking

type ValidationError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type InternalError struct {
	Message string
}

type ConflictError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e NotFoundError) Error() string { return e.Message }

func (e InternalError) Error() string { return e.Message }

func (e ConflictError) Error() string { return e.Message }
