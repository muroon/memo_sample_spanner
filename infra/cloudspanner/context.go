package cloudspanner

import (
	"context"

	"cloud.google.com/go/spanner"
)

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
	return context.WithValue(ctx, SpannerTransactionContext, spnCtx)
}

func setReadOnlyTransaction(
	ctx context.Context, tx *spanner.ReadOnlyTransaction,
) context.Context {
	spnCtx := &spannerTxContext{
		roTx: tx,
	}
	return context.WithValue(ctx, SpannerTransactionContext, spnCtx)
}

func setBatchReadOnlyTransaction(
	ctx context.Context, tx *spanner.BatchReadOnlyTransaction,
) context.Context {
	spnCtx := &spannerTxContext{
		broTx: tx,
	}
	return context.WithValue(ctx, SpannerTransactionContext, spnCtx)
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

	spnCtx, ok := ctx.Value(SpannerTransactionContext).(*spannerTxContext)
	if !ok {
		spnCtx = new(spannerTxContext)
	}

	return spnCtx
}
