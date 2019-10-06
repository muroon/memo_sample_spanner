package db

import (
	"context"
	"database/sql"
	"memo_sample_spanner/adapter/error"
	"memo_sample_spanner/infra/database"
	"memo_sample_spanner/infra/error"
)

var dbm *database.DBM

var errm apperror.ErrorManager

func init() {
	dbm = database.GetDBM()
	errm = apperrorsub.NewErrorManager()
}

// begin begin transaction
func begin(ctx context.Context) (context.Context, error) {
	return (*dbm).Begin(ctx)
}

// rollback rollback transaction
func rollback(ctx context.Context) (context.Context, error) {
	return (*dbm).Rollback(ctx)
}

// commit commit transaction
func commit(ctx context.Context) (context.Context, error) {
	return (*dbm).Commit(ctx)
}

// prepare prepare statement
func prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return (*dbm).Prepare(ctx, query)
}
