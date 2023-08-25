package custom_error

type CustomError struct {
	Field   string
	Message string
}

func (c CustomError) Error() string {
	return c.Message
}
