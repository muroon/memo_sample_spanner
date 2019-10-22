package cloudspanner

import (
	"context"
	"os"
)

var spnDB ISpannerDB

func init() {
	spnDB = newSpannerDB(
		setProjectID(os.Getenv("SPN_PROJECT_ID")),
		setInstanceID(os.Getenv("SPN_INSTANCE_ID")),
		setDatabaseID(os.Getenv("SPN_DATABASE_ID")),
	)
}

// OpenClient open spanner client
func OpenClient(ctx context.Context, local bool) error {
	if local {
		return spnDB.OpenClientLocal(ctx)
	}
	return spnDB.OpenClient(ctx)

}

// CloseClient close spanner client
func CloseClient() {
	spnDB.CloseClient()
}

func DB() ISpannerDB {
	return spnDB
}
