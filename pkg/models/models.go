package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("models: no requested snippet")

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
