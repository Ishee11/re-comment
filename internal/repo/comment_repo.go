package repo

import (
	"context"
	"github.com/evrone/go-clean-template/internal/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, c *entity.Comment) error
	Update(ctx context.Context, c *entity.Comment) error
	ListByEntity(ctx context.Context, entityID int64, page, limit int, sortAsc bool) (comments []*entity.Comment, total int64, err error)
	GetByID(ctx context.Context, id int64) (*entity.Comment, error)
}
