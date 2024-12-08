package mocks

import (
	"github.com/google/uuid"
	"github.com/joaooliveira247/go_olist_challenge/src/models"
	"github.com/lib/pq"
)

func NewMockBook() *models.Book {
	return &models.Book{
		Title:           "The Rust Programming Language",
		Edition:         1,
		PublicationYear: 2018,
	}
}

func NewMockBookOut() models.BookOut {
	return models.BookOut{
		Book: models.Book{
			ID:              uuid.New(),
			Title:           "the Rust Programming Language",
			Edition:         1,
			PublicationYear: 2018,
		},
		AuthorsName: pq.StringArray{"Carol Nichols", "Steve Klabnik"},
	}
}

func NewMockBooks() []models.BookOut {
	return []models.BookOut{
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "The Rust Programming Language",
				Edition:         uint8(1),
				PublicationYear: 2018,
			},
			AuthorsName: pq.StringArray{"Carol Nichols", "Steve Klabnik"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "The Go Programming Language",
				Edition:         uint8(2),
				PublicationYear: 2015,
			},
			AuthorsName: pq.StringArray{"Alan A. A. Donovan", "Brian W. Kernighan"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "Python Fluente",
				Edition:         uint8(1),
				PublicationYear: 2015,
			},
			AuthorsName: pq.StringArray{"Luciano Ramalho"},
		},
		{
			Book: models.Book{
				ID:              uuid.New(),
				Title:           "Python Fluente",
				Edition:         uint8(2),
				PublicationYear: 2022,
			},
			AuthorsName: pq.StringArray{"Luciano Ramalho"},
		},
	}
}
