package spanner

import (
	"context"
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/domain/repository"
	"memo_sample_spanner/infra/cloudspanner"
	"net/http"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
)

// NewMemoRepository get repository
func NewMemoRepository() repository.MemoRepository {
	return &memoRepository{}
}

// MemoRepository Memo's Repository Sub
type memoRepository struct {
}

// Save save Memo Data
func (m *memoRepository) Save(ctx context.Context, text string) (*model.Memo, error) {
	id, err := generateID()
	if err != nil {
		return nil, errm.Wrap(err, http.StatusInternalServerError)
	}

	mutations := make([]*spanner.Mutation, 0)

	memo := &model.Memo{
		MemoID: id,
		Text:   spanner.NullString{StringVal: text, Valid: true},
	}

	mutations = append(mutations, memo.Insert(ctx))

	return memo, cloudspanner.DB().Apply(ctx, mutations)
}

// Get get Memo Data by ID
func (m memoRepository) Get(ctx context.Context, id string) (*model.Memo, error) {
	return model.FindMemo(ctx, yoRODB(), id)
}

// GetAll get all Memo Data
func (m *memoRepository) GetAll(ctx context.Context) ([]*model.Memo, error) {
	list := make([]*model.Memo, 0)

	yoDB := yoRODB()

	iter := yoDB.Read(ctx, "Memo", spanner.AllKeys(), model.MemoColumns())
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		memo := new(model.Memo)

		if err := row.Columns(&memo.MemoID, &memo.Text); err != nil {
			return list, err
		}

		list = append(list, memo)
	}

	return list, nil
}

// Search search memo by text
func (m *memoRepository) Search(ctx context.Context, text string) ([]*model.Memo, error) {
	list := make([]*model.Memo, 0)

	yoDB := yoRODB()

	stmt := spanner.Statement{
		SQL: `SELECT Memo.* FROM Memo WHERE STARTS_WITH(Text, @text)`,
		Params: map[string]interface{}{
			"text": text,
		},
	}

	iter := yoDB.Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		memo := new(model.Memo)

		if err := row.Columns(&memo.MemoID, &memo.Text); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, memo)
	}

	return list, nil
}

// GetAllByIDs get all Memo Data by ID
func (m *memoRepository) GetAllByIDs(ctx context.Context, ids []string) ([]*model.Memo, error) {
	list := make([]*model.Memo, 0, len(ids))

	yoDB := yoRODB()

	keys := make([]spanner.KeySet, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, spanner.Key{id})
	}
	keySet := spanner.KeySets(keys...)

	iter := yoDB.Read(ctx, "Memo", keySet, model.MemoColumns())
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		memo := new(model.Memo)

		if err := row.Columns(&memo.MemoID, &memo.Text); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, memo)
	}

	return list, nil
}
