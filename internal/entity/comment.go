package entity

import (
	"errors"
	"strings"
	"time"
)

// ErrEmptyCommentText и ErrNotAuthor ошибки пустого комментария и неверного автора
var ErrEmptyCommentText = errors.New("text is empty")
var ErrNotAuthor = errors.New("not author")

// Comment представляет комментарий пользователя к какой-либо сущности.
// Внутри хранит ID комментария, ID пользователя, текст, дату создания и изменения.
type Comment struct {
	ID        int64
	EntityID  int64
	UserID    int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment валидация и создание комментария
func NewComment(entityID, userID int64, text string) (*Comment, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrEmptyCommentText
	}

	now := time.Now().UTC()

	return &Comment{
		EntityID:  entityID,
		UserID:    userID,
		Text:      text,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateCommentText валидация, изменения текста и даты обновления комментария
func (c *Comment) UpdateCommentText(userID int64, text string) error {
	if userID != c.UserID {
		return ErrNotAuthor
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return ErrEmptyCommentText
	}
	c.Text = text
	c.UpdatedAt = time.Now().UTC()
	return nil
}
