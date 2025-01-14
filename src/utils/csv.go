package utils

import (
	"encoding/csv"
	"os"

	"github.com/joaooliveira247/go_olist_challenge/src/models"
)

func ParseAuthorsFromCSV(path string, header bool) ([]models.Author, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	if header {
		lines = lines[1:]
	}

	var authors []models.Author

	for _, line := range lines {
		author := models.Author{Name: line[0]}
		authors = append(authors, author)
	}

	return authors, nil
}
