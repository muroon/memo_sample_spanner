package db

import (
	"context"
	"database/sql"
	"memo_sample/domain/model"
	"memo_sample/domain/repository"
	"strconv"
	"strings"
)

// NewMemoRepository get repository
func NewMemoRepository() repository.MemoRepository {
	return memoRepository{}
}

// memoRepository Memo's Repository Sub
type memoRepository struct{}

// Save save Memo Data
func (m memoRepository) Save(ctx context.Context, text string) (*model.Memo, error) {
	var err error
	var res sql.Result
	query := "insert into memo(text) values(?)"
	stmt, err := prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err = stmt.ExecContext(ctx, text)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.Get(ctx, int(id))
}

// Get get Memo Data by ID
func (m memoRepository) Get(ctx context.Context, id int) (*model.Memo, error) {
	mem := &model.Memo{}
	var err error
	query := "select * from memo where id = ?"
	stmt, err := prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&mem.ID, &mem.Text)
	if err != nil {
		return nil, err
	}

	return mem, err
}

// GetAll get all Memo Data
func (m memoRepository) GetAll(ctx context.Context) ([]*model.Memo, error) {
	var rows *sql.Rows
	var err error
	query := "select * from memo"
	stmt, err := prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err = stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	return m.getModelList(rows)
}

// Search search memo by text
func (m memoRepository) Search(ctx context.Context, text string) ([]*model.Memo, error) {
	var rows *sql.Rows
	var err error
	query := "select * from memo where text like '%" + text + "%'"
	stmt, err := prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err = stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	return m.getModelList(rows)
}

// GetAllByIDs get all Memo Data by ID
func (m memoRepository) GetAllByIDs(ctx context.Context, ids []int) ([]*model.Memo, error) {
	idvs := []string{}
	for _, id := range ids {
		idvs = append(idvs, strconv.Itoa(id))
	}

	query := "select * from memo where id in (" + strings.Join(idvs, ",") + ")"

	var rows *sql.Rows
	var err error
	stmt, err := prepare(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err = stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	return m.getModelList(rows)
}

// getModelList get model list
func (m memoRepository) getModelList(rows *sql.Rows) ([]*model.Memo, error) {
	list := []*model.Memo{}
	for rows.Next() {
		mem := &model.Memo{}
		err := rows.Scan(&mem.ID, &mem.Text)
		if err != nil {
			return list, err
		}
		list = append(list, mem)
	}

	return list, nil
}
