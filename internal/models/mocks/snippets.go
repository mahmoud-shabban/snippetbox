package mocks

import (
	"time"

	"github.com/mahmoud-shabban/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "Mock snippet title",
	Content: "mock snippet content",
	Created: time.Now(),
	Expires: time.Now().Add(time.Hour),
}

type SnippetModel struct{}

func (s *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}
func (s *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, nil
	}
}

func (s *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
