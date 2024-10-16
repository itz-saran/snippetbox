package models

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Snippet struct {
	ID          int
	Title       string
	Content     string
	ContentHTML template.HTML
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expiresIn int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING id`
	createdAt := pq.FormatTimestamp(time.Now())
	expiresAt := pq.FormatTimestamp(time.Now().AddDate(0, 0, expiresIn))
	var result int
	err := m.DB.QueryRow(stmt, title, content, createdAt, expiresAt).Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := "SELECT id, title, content, created_at, expires_at FROM snippets WHERE expires_at > $1 AND id=$2"
	s := &Snippet{}
	err := m.DB.QueryRow(stmt, pq.FormatTimestamp(time.Now()), id).Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// TODO use functions instead of extra variable
	s.ContentHTML = template.HTML(strings.ReplaceAll(s.Content, "\\n", "<br>"))
	fmt.Println(s.ContentHTML)
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created_at, expires_at FROM snippets WHERE expires_at > $1 ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, pq.FormatTimestamp(time.Now()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
