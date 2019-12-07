package cloudspanner

import (
	"context"
	"memo_sample_spanner/domain/app"
	"time"

	"cloud.google.com/go/spanner"
)

func (d *spannerDB) ReadWriteTransaction(
	ctx context.Context,
	f func(ctx context.Context) error,
) (time.Time, error) {
	return d.client.ReadWriteTransaction(
		ctx,
		func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
			ctx = setReadWriteTransaction(ctx, tx)
			return f(ctx)
		})
}

func (d *spannerDB) ReadOnlyTransaction(
	ctx context.Context,
	f func(ctx context.Context) error,
) error {
	tx := d.client.ReadOnlyTransaction()
	ctx = setReadOnlyTransaction(ctx, tx)
	defer tx.Close()

	return f(ctx)
}

func (d *spannerDB) BatchReadOnlyTransaction(
	ctx context.Context,
	f func(ctx context.Context) error,
) error {
	tx, err := d.client.BatchReadOnlyTransaction(ctx, spanner.StrongRead()) // TODO:
	if err != nil {
		return d.errManager.Wrap(err, app.DBError)
	}
	defer tx.Close()
	ctx = setBatchReadOnlyTransaction(ctx, tx)

	return f(ctx)
}
