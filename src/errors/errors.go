package errors

import "fmt"

type BaseError struct {
	Msg      string
	Resource string
}

func (error *BaseError) Error() string {
	return fmt.Sprintf("%s %s", error.Resource, error.Msg)
}

type GenericError struct {
	BaseError
}
