package spanner

import (
	"context"
	"memo_sample_spanner/infra/cloudspanner"
)

// connectTestDB DB接続
func connectTestDB(ctx context.Context) {
	_ = cloudspanner.OpenClient(ctx, false)
}

// closeTestDB DB切断
func closeTestDB() {
	cloudspanner.CloseClient()
}
