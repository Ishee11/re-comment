package repo

import (
	"github.com/evrone/go-clean-template/internal/entity"
)

type CommentRepository interface {
	CreateComment(c *entity.Comment) error
	UpdateComment(c *entity.Comment) error
	ListCommentByEntity(entityID int64, page int64, limit int64, sortAsc bool) ([]*entity.Comment, error)
	GetCommentByID(id int64) (*entity.Comment, error)
}
