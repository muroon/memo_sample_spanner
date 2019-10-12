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

// NewTagRepository get repository
func NewTagRepository() repository.TagRepository {
	return &tagRepository{}
}

// tagRepository Tag's Repository Sub
type tagRepository struct{}

// Save save Tag Data
func (m *tagRepository) Save(ctx context.Context, title string) (*model.Tag, error) {
	id, err := generateID()
	if err != nil {
		return nil, errm.Wrap(err, http.StatusInternalServerError)
	}

	mutations := make([]*spanner.Mutation, 0)

	memo := &model.Tag{
		TagID: id,
		Title: spanner.NullString{StringVal: title, Valid: true},
	}

	mutations = append(mutations, memo.Insert(ctx))

	return memo, cloudspanner.DB().Apply(ctx, mutations)
}

// Get get Tag Data by ID
func (m tagRepository) Get(ctx context.Context, id string) (*model.Tag, error) {
	return model.FindTag(ctx, yoRODB(), id)
}

// GetAll get all Tag Data
func (m *tagRepository) GetAll(ctx context.Context) ([]*model.Tag, error) {
	list := make([]*model.Tag, 0)

	yoDB := yoRODB()

	iter := yoDB.Read(ctx, "Tag", spanner.AllKeys(), model.TagColumns())
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		tag := new(model.Tag)

		if err := row.Columns(&tag.TagID, &tag.Title); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, tag)
	}

	return list, nil
}

// Search search tag by text
func (m *tagRepository) Search(ctx context.Context, title string) ([]*model.Tag, error) {
	list := make([]*model.Tag, 0)

	yoDB := yoRODB()

	stmt := spanner.Statement{
		SQL: `SELECT Tag.* FROM Tag WHERE STARTS_WITH(Title, @title)`,
		Params: map[string]interface{}{
			"title": title,
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

		tag := new(model.Tag)

		if err := row.Columns(&tag.TagID, &tag.Title); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, tag)
	}

	return list, nil
}

// SaveTagAndTag save tag and memo link
func (m *tagRepository) SaveTagAndMemo(ctx context.Context, tagID, memoID string) error {
	mutations := make([]*spanner.Mutation, 0)

	tm := &model.TagMemo{
		TagID:  tagID,
		MemoID: memoID,
	}

	mutations = append(mutations, tm.Insert(ctx))

	return cloudspanner.DB().Apply(ctx, mutations)
}

// GetAllByMemoID get all Tag Data By MemoID
func (m *tagRepository) GetAllByMemoID(ctx context.Context, memoID string) ([]*model.Tag, error) {
	list := make([]*model.Tag, 0)

	yoDB := yoRODB()

	stmt := spanner.Statement{
		SQL: `SELECT Tag.* FROM Tag INNER JOIN TagMemo ON Tag.TagID = TagMemo.TagID WHERE TagMemo.MemoID = @memoID`,
		Params: map[string]interface{}{
			"memoID": memoID,
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

		tag := new(model.Tag)

		if err := row.Columns(&tag.TagID, &tag.Title); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, tag)
	}

	return list, nil
}

// SearchMemoIDsByTitle search memo ids by tag's title
func (m *tagRepository) SearchMemoIDsByTitle(ctx context.Context, title string) ([]string, error) {
	list := make([]string, 0)

	yoDB := yoRODB()

	stmt := spanner.Statement{
		SQL: `SELECT TagMemo.MemoID FROM Tag INNER JOIN TagMemo ON Tag.TagID = TagMemo.TagID WHERE STARTS_WITH(Tag.Title, @title)`,
		Params: map[string]interface{}{
			"title": title,
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

		var memoID string

		if err := row.Columns(&memoID); err != nil {
			return list, errm.Wrap(err, http.StatusInternalServerError)
		}

		list = append(list, memoID)
	}

	return list, nil
}

// SearchMemoAndTagByTagTitle search memo ids by tag's title
func (m *tagRepository) SearchMemoAndTagByTagTitle(ctx context.Context, title string) (
	[]*model.Memo, []*model.Tag, error,
) {
	memos := make([]*model.Memo, 0)
	tags := make([]*model.Tag, 0)

	yoDB := yoRODB()

	stmt := spanner.Statement{
		SQL: `SELECT Memo.*, Tag.* FROM Memo, Tag, TagMemo WHERE Memo.MemoID = TagMemo.MemoID AND Tag.TagID = TagMemo.TagID AND STARTS_WITH(Tag.Title, @title)`,
		Params: map[string]interface{}{
			"title": title,
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
			return memos, tags, errm.Wrap(err, http.StatusInternalServerError)
		}

		memo := new(model.Memo)
		tag := new(model.Tag)

		if err := row.Columns(&memo.MemoID, &memo.Text, &tag.TagID, &tag.Title); err != nil {
			return memos, tags, errm.Wrap(err, http.StatusInternalServerError)
		}

		memos = append(memos, memo)
		tags = append(tags, tag)
	}

	return memos, tags, nil
}
