package dto

import (
	"fmt"
	"strings"
)

type AuthorQueryParams struct {
	ID   string `form:"authorID"`
	Name string `form:"name"`
}

type BookQueryParams struct {
	AuthorID        string `form:"authorID,omitempty"`
	BookID          string `form:"bookID,omitempty"`
	Title           string `form:"title,omitempty"`
	Edition         uint8  `form:"edition,omitempty"`
	PublicationYear uint   `form:"publicationYear,omitempty"`
}

func (query *BookQueryParams) AsQuery() string {
	whereClauses := []string{}

	if query.Title != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.title = '%s'`, query.Title))
	}
	if query.Edition != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.edition = %d`, query.Edition))
	}
	if query.PublicationYear != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf(`b.publication_year = %d`, query.PublicationYear))
	}

	return strings.Join(whereClauses, " AND ")
}

func (query *BookQueryParams) IsEmpty() bool {
	return query.Title == "" && query.Edition == 0 && query.PublicationYear == 0
}
