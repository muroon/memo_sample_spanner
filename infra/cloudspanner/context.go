package cloudspanner

import (
	"context"
	"memo_sample_spanner/domain/app"

	"cloud.google.com/go/spanner"
)

// SpannerTransactionContext context of spanner transaction
const SpannerTransactionContext string = "spnContext"

type spannerTxContext struct {
	rwTx  *spanner.ReadWriteTransaction
	roTx  *spanner.ReadOnlyTransaction
	broTx *spanner.BatchReadOnlyTransaction
}

func setReadWriteTransaction(
	ctx context.Context, tx *spanner.ReadWriteTransaction,
) context.Context {
	spnCtx := &spannerTxContext{
		rwTx: tx,
	}
	return context.WithValue(ctx, app.ContextKey(app.SpannerTransaction), spnCtx)
}

func setReadOnlyTransaction(
	ctx context.Context, tx *spanner.ReadOnlyTransaction,
) context.Context {
	spnCtx := &spannerTxContext{
		roTx: tx,
	}
	return context.WithValue(ctx, app.ContextKey(app.SpannerTransaction), spnCtx)
}

func setBatchReadOnlyTransaction(
	ctx context.Context, tx *spanner.BatchReadOnlyTransaction,
) context.Context {
	spnCtx := &spannerTxContext{
		broTx: tx,
	}
	return context.WithValue(ctx, app.ContextKey(app.SpannerTransaction), spnCtx)
}

func getAllTransactions(
	ctx context.Context,
) (*spanner.ReadWriteTransaction,
	*spanner.ReadOnlyTransaction,
	*spanner.BatchReadOnlyTransaction,
) {
	spnCtx := getSpannerContext(ctx)
	return spnCtx.rwTx, spnCtx.roTx, spnCtx.broTx
}

func getSpannerContext(
	ctx context.Context,
) *spannerTxContext {

	spnCtx, ok := ctx.Value(app.ContextKey(app.SpannerTransaction)).(*spannerTxContext)
	if !ok {
		spnCtx = new(spannerTxContext)
	}

	return spnCtx
}
