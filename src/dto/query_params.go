package dto

import (
	"fmt"
	"strings"
)

type AuthorQueryParam struct {
	Name string `form:"name" binding:"min=2"`
}

type BookQueryParam struct {
	Title           string `form:"title,omitempty"`
	Edition         uint8  `form:"edition,omitempty"`
	PublicationYear uint   `form:"publicationYear,omitempty"`
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

func (query *BookQueryParam) IsEmpty() bool {
	return query.Title == "" && query.Edition == 0 && query.PublicationYear == 0
}
