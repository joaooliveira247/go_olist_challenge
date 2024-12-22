package dto

import (
	"fmt"
	"strings"
)

type AuthorQueryParam struct {
	Name string `form:"name" binding:"min=2"`
}

type BookQueryParam struct {
	Title           string `form:"title,omitempty" json:"title,omitempty"`
	Edition         uint8  `form:"edition,omitempty" json:"edition,omitempty"`
	PublicationYear uint   `form:"publicationYear,omitempty" json:"publication_year,omitempty"`
}

func (query *BookQueryParam) AsQuery() string {
	whereClauses := []string{}

	if query.Title != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.title = %s`, query.Title))
	}
	if query.Edition != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.edition = %d`, query.Edition))
	}
	if query.PublicationYear != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.publication_year = %d`, query.PublicationYear))
	}

	return strings.Join(whereClauses, " AND ")
}
