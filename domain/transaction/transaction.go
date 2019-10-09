package transaction

import (
	"context"
	"time"
)

type ITransaction interface {
	ReadWriteTransaction(
		ctx context.Context,
		f func(ctx context.Context) error,
	) (time.Time, error)

	ReadOnlyTransaction(
		ctx context.Context,
		f func(ctx context.Context) error,
	) error

	BatchReadOnlyTransaction(
		ctx context.Context,
		f func(ctx context.Context) error,
	) error
}
