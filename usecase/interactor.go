package usecase

import (
	"context"
	"memo_sample_spanner/domain/transaction"
	"memo_sample_spanner/usecase/input"
)

// IInteractor interface of interacor
type IInteractor interface {
	PostMemo(ctx context.Context, ipt input.PostMemo)
	GetMemos(ctx context.Context)
	PostMemoAndTags(ctx context.Context, ipt input.PostMemoAndTags)
	SearchTagsAndMemos(ctx context.Context, ipt input.SearchTagsAndMemos)
}

// NewInteractor new Interactor
func NewInteractor(
	pre Presenter,
	tx transaction.ITransaction,
	memo Memo,
) IInteractor {
	return &interactor{
		pre:  pre,
		tx:   tx,
		memo: memo,
	}
}

// interactor usecase interactor
type interactor struct {
	pre  Presenter
	tx   transaction.ITransaction
	memo Memo
}

// PostMemo post memo
func (i *interactor) PostMemo(ctx context.Context, ipt input.PostMemo) {
	var id string
	_, err := i.tx.ReadWriteTransaction(ctx,
		func(ctx context.Context) error {
			err := i.memo.ValidatePost(ipt)
			if err != nil {
				return err
			}

			id, err = i.memo.Post(ctx, ipt)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		i.pre.ViewError(ctx, err)
	}

	iptf := &input.GetMemo{ID: id}
	memo, err := i.memo.GetMemo(ctx, *iptf)
	if err != nil {
		i.pre.ViewError(ctx, err)
		return
	}

	i.pre.ViewMemo(ctx, memo)
}

// GetMemos get all memos
func (i *interactor) GetMemos(ctx context.Context) {

	memos, err := i.memo.GetAllMemoList(ctx)
	if err != nil {
		i.pre.ViewError(ctx, err)
		return
	}

	i.pre.ViewMemoList(ctx, memos)
}

// PostMemoAndTags save memo and tags
func (i *interactor) PostMemoAndTags(ctx context.Context, ipt input.PostMemoAndTags) {
	_, err := i.tx.ReadWriteTransaction(ctx,
		func(ctx context.Context) error {

			err := i.memo.ValidatePostMemoAndTags(ipt)
			if err != nil {
				return err
			}

			memo, tags, err := i.memo.PostMemoAndTags(ctx, ipt)
			if err != nil {
				return err
			}

			i.pre.ViewPostMemoAndTagsResult(ctx, memo, tags)
			return nil
		},
	)

	if err != nil {
		i.pre.ViewError(ctx, err)
	}
}

// SearchTagsAndMemos save memo and tags
func (i *interactor) SearchTagsAndMemos(ctx context.Context, ipt input.SearchTagsAndMemos) {
	err := i.tx.ReadOnlyTransaction(ctx, func(ctx context.Context) error {
		memos, tags, err := i.memo.SearchTagsAndMemos(ctx, ipt)
		if err != nil {
			return err
		}

		i.pre.ViewSearchTagsAndMemosResult(ctx, memos, tags)
		return nil
	})

	if err != nil {
		i.pre.ViewError(ctx, err)
	}
}
