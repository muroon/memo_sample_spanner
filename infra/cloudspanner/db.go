package cloudspanner

import (
	"context"
	"fmt"
	apperror "memo_sample_spanner/infra/error"
	"time"

	"cloud.google.com/go/spanner"

	"github.com/gcpug/handy-spanner/fake"
)

// ISpannerDB spannerDB interface
type ISpannerDB interface {
	OpenClient(ctx context.Context) error
	OpenClientLocal(ctx context.Context) error
	CloseClient()
	Client() *spanner.Client

	Apply(
		ctx context.Context, mutations []*spanner.Mutation,
	) error

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

// NewSpannerDB new spanner DB interface
func newSpannerDB(opts ...Option) ISpannerDB {
	s := new(spannerDB)
	for _, opt := range opts {
		opt(s)
	}
	s.errManager = apperror.NewErrorManager()
	return s
}

// Option new spanner DB option
type Option func(*spannerDB)

// setProjectID projectID option
func setProjectID(id string) Option {
	return func(d *spannerDB) {
		d.projectID = id
	}
}

// setInstanceID instanceID option
func setInstanceID(id string) Option {
	return func(d *spannerDB) {
		d.instanceID = id
	}
}

// setDatabaseID databaseID option
func setDatabaseID(id string) Option {
	return func(d *spannerDB) {
		d.databaseID = id
	}
}

type spannerDB struct {
	projectID   string
	instanceID  string
	databaseID  string
	client      *spanner.Client
	errManager  apperror.ErrorManager
	localServer *fake.Server
}

func (d *spannerDB) OpenClient(ctx context.Context) error {
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		d.projectID, d.instanceID, d.databaseID,
	)

	var err error
	d.client, err = spanner.NewClient(ctx, databaseName)
	return err
}

func (d *spannerDB) CloseClient() {
	if d.localServer != nil {
		d.localServer.Stop()
	}
	if d.client != nil {
		d.client.Close()
	}
}

func (d *spannerDB) Client() *spanner.Client {
	return d.client
}

func (d *spannerDB) ReadRow(
	ctx context.Context, table string, key spanner.Key, columns []string,
) (*spanner.Row, error) {
	rwTx, roTx, broTx := getAllTransactions(ctx)
	switch {
	case rwTx != nil:
		return rwTx.ReadRow(ctx, table, key, columns)
	case roTx != nil:
		return roTx.ReadRow(ctx, table, key, columns)
	case broTx != nil:
		return broTx.ReadRow(ctx, table, key, columns)
	}
	return d.client.Single().ReadRow(ctx, table, key, columns)
}

func (d *spannerDB) Read(
	ctx context.Context, table string, keys spanner.KeySet, columns []string,
) *spanner.RowIterator {
	rwTx, roTx, broTx := getAllTransactions(ctx)
	switch {
	case rwTx != nil:
		return rwTx.Read(ctx, table, keys, columns)
	case roTx != nil:
		return roTx.Read(ctx, table, keys, columns)
	case broTx != nil:
		return broTx.Read(ctx, table, keys, columns)
	}

	iter := d.client.Single().Read(ctx, table, keys, columns)
	return iter
}

func (d *spannerDB) ReadUsingIndex(
	ctx context.Context, table, index string, keys spanner.KeySet, columns []string,
) (ri *spanner.RowIterator) {
	rwTx, roTx, broTx := getAllTransactions(ctx)
	switch {
	case rwTx != nil:
		return rwTx.ReadUsingIndex(ctx, table, index, keys, columns)
	case roTx != nil:
		return roTx.ReadUsingIndex(ctx, table, index, keys, columns)
	case broTx != nil:
		return broTx.ReadUsingIndex(ctx, table, index, keys, columns)
	}

	iter := d.client.Single().ReadUsingIndex(ctx, table, index, keys, columns)
	return iter
}

func (d *spannerDB) Query(
	ctx context.Context, statement spanner.Statement,
) *spanner.RowIterator {
	rwTx, roTx, broTx := getAllTransactions(ctx)
	switch {
	case rwTx != nil:
		return rwTx.Query(ctx, statement)
	case roTx != nil:
		return roTx.Query(ctx, statement)
	case broTx != nil:
		return broTx.Query(ctx, statement)
	}

	return d.client.Single().Query(ctx, statement)
}

func (d *spannerDB) Apply(
	ctx context.Context, mutations []*spanner.Mutation,
) error {
	rwTx, _, _ := getAllTransactions(ctx)
	switch {
	case rwTx != nil:
		return rwTx.BufferWrite(mutations)
	}
	_, err := d.client.Apply(ctx, mutations)
	return err
}
