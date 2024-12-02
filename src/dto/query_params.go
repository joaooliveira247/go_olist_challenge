package dto

type AuthorQueryParam struct {
	Name string `form:"name" binding:"min=2"`
}

type BookQueryParam struct {
	Title           string `form:"title,omitempty" json:"title,omitempty"`
	Edition         uint8  `form:"edition,omitempty" json:"edition,omitempty"`
	PublicationYear uint   `form:"publicationYear,omitempty" json:"publication_year,omitempty"`
	Author          string `form:"author,omitempty" json:"author,omitempty"`
}
