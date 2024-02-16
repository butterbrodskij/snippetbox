package postgres

import (
	"database/sql"
	"errors"

	"mysnippetbox.com/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	toExec := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1, $2, NOW(), NOW() + make_interval(days => $3)) RETURNING id;`

	var id int
	err := m.DB.QueryRow(toExec, title, content, expires).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	toExec := "SELECT * FROM snippets WHERE expires > NOW() AND id = $1"

	snip := new(models.Snippet)
	err := m.DB.QueryRow(toExec, id).Scan(&snip.Id, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return snip, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	toExec := "SELECT * FROM snippets WHERE expires > NOW() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(toExec)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var snippets []*models.Snippet

	for rows.Next() {
		snip := new(models.Snippet)
		err = rows.Scan(&snip.Id, &snip.Title, &snip.Content, &snip.Created, &snip.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
