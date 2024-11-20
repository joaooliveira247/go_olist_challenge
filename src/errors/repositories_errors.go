package errors

type AlreadyExists struct {
	BaseError
}

type NotFound struct {
	BaseError
}

var (
	AuthorAlreadyExists       = AlreadyExists{BaseError{"author", "already exists"}}
	AuthorGenericError        = GenericError{BaseError{"author", "generic error"}}
	BookAuthorGenericError    = GenericError{BaseError{"book_author", "generic error"}}
	AuthorNotFound            = NotFound{BaseError{"author", "not found"}}
	RelationshipAlreadyExists = AlreadyExists{BaseError{"relationship", "already exists"}}
)
