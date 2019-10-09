package repository

import (
	"context"
	"memo_sample_spanner/domain/model"
)

// TagRepository Tag's Repository
type TagRepository interface {
	Save(ctx context.Context, title string) (*model.Tag, error)
	Get(ctx context.Context, id string) (*model.Tag, error)
	GetAll(ctx context.Context) ([]*model.Tag, error)
	Search(ctx context.Context, title string) ([]*model.Tag, error)
	SaveTagAndMemo(ctx context.Context, tagID, memoID string) error
	GetAllByMemoID(ctx context.Context, id string) ([]*model.Tag, error)
	SearchMemoIDsByTitle(ctx context.Context, title string) ([]string, error)
}
