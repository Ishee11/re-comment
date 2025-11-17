package comment

import (
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

type UseCase struct {
	repo repo.CommentRepository
}

func (uc *UseCase) CreateComment(userID, entityID int64, text string) (*entity.Comment, error) {
	c, err := entity.NewComment(entityID, userID, text)
	if err != nil {
		return nil, err
	}
	if err = uc.repo.CreateComment(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *UseCase) UpdateComment(id, userID int64, text string) error {
	c, err := uc.repo.GetCommentByID(id)
	if err != nil {
		return err
	}
	if err = c.UpdateCommentText(userID, text); err != nil {
		return err
	}
	return uc.repo.UpdateComment(c)
}

func (uc *UseCase) ListComments(entityID, page, limit int64, sortAsc bool) ([]*entity.Comment, error) {
	return uc.repo.ListCommentByEntity(entityID, page, limit, sortAsc)
}
