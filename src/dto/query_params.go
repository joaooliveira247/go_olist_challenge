package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AuthorQueryParam struct {
	Name string `form:"name" binding:"min=2"`
}

type BookQueryParam struct {
	Title           string    `form:"title,omitempty" json:"title,omitempty"`
	Edition         uint8     `form:"edition,omitempty" json:"edition,omitempty"`
	PublicationYear uint      `form:"publicationYear,omitempty" json:"publication_year,omitempty"`
}

func (query *BookQueryParam) AsQuery() (map[string]interface{}, error) {
	bjson, err := json.Marshal(query)

	if err != nil {
		return nil, err
	}

	var queries map[string]interface{}

	if err := json.Unmarshal(bjson, &queries); err != nil {
		return nil, err
	}

	return queries, nil
}
