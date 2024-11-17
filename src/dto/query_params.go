package dto

type AuthorQueryParam struct {
	Name string `form:"name" binding:"min=2"`
}
