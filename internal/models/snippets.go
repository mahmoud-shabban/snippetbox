package models

import (
	"database/sql"
	"errors"
	"time"
)

// snippet data model
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// wraper around db connection pool
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {

	// create sql statement with placeholder for safer execution
	stmt := `
			INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, Now(), DATE_ADD(NOW(), INTERVAL ? DAY)
			`
	result, err := m.DB.Exec(stmt, title, content, expires) // result implements result interface from database/sql pkg
	if err != nil {
		return 0, err
	}

	// id is the id of the inserted item
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `
			SELECT id, title, content, created, expires
			FROM snippets
			WHERE expires > NOW() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	snippet := Snippet{}
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return snippet, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `
			SELECT id, title, content, created, expires 
			FROM snippets
			WHERE expires > NOW()
			ORDER BY created DESC
			LIMIT 10
			`
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	snippets := make([]Snippet, 0)

	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
