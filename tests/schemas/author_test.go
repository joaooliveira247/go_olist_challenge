package schemas_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/joaooliveira247/go_olist_challenge/src/schemas"
	"github.com/stretchr/testify/assert"
)

var validate = validator.New()

func TestAuthorInSchemaValid(t *testing.T) {
	author := schemas.AuthorIn{
		Name: "Luciano Ramalho",
	}

	err := validate.Struct(author)

	assert.Nil(t, err)
}
