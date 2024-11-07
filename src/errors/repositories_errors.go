package errors

type AlreadyExists struct {
	BaseError
}

var AuthorAlreadyExists = AlreadyExists{BaseError{"author", "already exists"}}
