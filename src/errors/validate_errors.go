package errors

type FieldInvalid struct {
	BaseError
}

var (
	FieldNameInvalid = FieldInvalid{BaseError{"field", "name is not valid in body"}}
)
