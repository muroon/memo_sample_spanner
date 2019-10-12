package usecase

import (
	"context"
	"fmt"
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/domain/repository"
	"memo_sample_spanner/infra/error"
	"memo_sample_spanner/usecase/input"
	"net/http"
)

// Memo memo related interface
type Memo interface {
	ValidatePost(ipt input.PostMemo) error
	Post(ctx context.Context, ipt input.PostMemo) (string, error)
	ValidateGet(ipt input.GetMemo) error
	GetMemo(ctx context.Context, ipt input.GetMemo) (*model.Memo, error)
	GetAllMemoList(ctx context.Context) ([]*model.Memo, error)
	ValidatePostMemoAndTags(ipt input.PostMemoAndTags) error
	PostMemoAndTags(ctx context.Context, ipt input.PostMemoAndTags) (*model.Memo, []*model.Tag, error)
	GetTagsByMemo(ctx context.Context, ipt input.GetTagsByMemo) ([]*model.Tag, error)
	SearchTagsAndMemos(ctx context.Context, ipt input.SearchTagsAndMemos) ([]*model.Memo, []*model.Tag, error)
}

// NewMemo generate memo instance
func NewMemo(
	memoRepository repository.MemoRepository,
	tagRepository repository.TagRepository,
	errm apperror.ErrorManager,
) Memo {
	return memo{
		memoRepository,
		tagRepository,
		errm,
	}
}

type memo struct {
	memoRepository repository.MemoRepository
	tagRepository  repository.TagRepository
	errm           apperror.ErrorManager
}

func (m memo) ValidatePost(ipt input.PostMemo) error {
	if ipt.Text == "" {
		err := fmt.Errorf("text parameter is invalid. %s", ipt.Text)
		return m.errm.Wrap(
			err,
			http.StatusBadRequest,
		)
	}

	return nil
}

func (m memo) Post(ctx context.Context, ipt input.PostMemo) (string, error) {
	mo, err := m.memoRepository.Save(ctx, ipt.Text)
	if err != nil {
		return "", err
	}
	return mo.MemoID, err
}

func (m memo) ValidateGet(ipt input.GetMemo) error {
	if ipt.ID == "" {
		err := fmt.Errorf("ID parameter is invalid. %s", ipt.ID)
		return m.errm.Wrap(
			err,
			http.StatusBadRequest,
		)
	}

	return nil
}

func (m memo) GetMemo(ctx context.Context, ipt input.GetMemo) (*model.Memo, error) {
	me, err := m.memoRepository.Get(ctx, ipt.ID)
	if err != nil {
		return nil, err
	}

	return me, nil
}

func (m memo) GetAllMemoList(ctx context.Context) ([]*model.Memo, error) {
	list, err := m.memoRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m memo) ValidatePostMemoAndTags(ipt input.PostMemoAndTags) error {
	if ipt.MemoText == "" {
		err := fmt.Errorf("text parameter(MemoText) is invalid. %s", ipt.MemoText)
		return m.errm.Wrap(
			err,
			http.StatusBadRequest,
		)
	}

	for _, title := range ipt.TagTitles {
		if title == "" {
			err := fmt.Errorf("text parameter(TagTitles) is invalid. %s", title)
			return m.errm.Wrap(
				err,
				http.StatusBadRequest,
			)
		}
	}

	return nil
}

func (m memo) PostMemoAndTags(ctx context.Context, ipt input.PostMemoAndTags) (*model.Memo, []*model.Tag, error) {
	tags := make([]*model.Tag, 0)

	// Memo
	mo, err := m.memoRepository.Save(ctx, ipt.MemoText)
	if err != nil {
		return nil, nil, err
	}

	for _, title := range ipt.TagTitles {
		// Tag
		tg, err := m.tagRepository.Save(ctx, title)
		if err != nil {
			return nil, nil, err
		}
		tags = append(tags, tg)

		// MemoTag
		err = m.tagRepository.SaveTagAndMemo(ctx, tg.TagID, mo.MemoID)
		if err != nil {
			return nil, nil, err
		}
	}

	return mo, tags, nil
}

func (m memo) GetTagsByMemo(ctx context.Context, ipt input.GetTagsByMemo) ([]*model.Tag, error) {
	return m.tagRepository.GetAllByMemoID(ctx, ipt.ID)
}

func (m memo) SearchTagsAndMemos(ctx context.Context, ipt input.SearchTagsAndMemos) ([]*model.Memo, []*model.Tag, error) {
	return m.tagRepository.SearchMemoAndTagByTagTitle(ctx, ipt.TagTitle)
}
