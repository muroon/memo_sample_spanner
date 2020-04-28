package cloudspanner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/spanner"
	"github.com/gcpug/handy-spanner/fake"
	"google.golang.org/api/option"
)

const shemaFile = "migration/schema.sql"

func (d *spannerDB) OpenClientLocal(ctx context.Context) error {
	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		d.projectID, d.instanceID, d.databaseID,
	)

	srv, conn, err := fake.Run()
	if err != nil {
		return err
	}
	d.localServer = srv

	d.client, err = spanner.NewClient(ctx, databaseName, option.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	return d.initalizeTables(ctx, databaseName)
}

func (d *spannerDB) initalizeTables(ctx context.Context, databaseName string) error {
	schema, err := readSchema(ctx)
	if err != nil {
		return err
	}

	return d.localServer.ParseAndApplyDDL(ctx, databaseName, strings.NewReader(schema))
	return nil
}

func readSchema(ctx context.Context) (string, error) {
	var schema string
	f, err := os.Open(shemaFile)
	if err != nil {
		return schema, err
	}
	defer func() {
		_ = f.Close()
	}()

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return schema, err
	}
	schema = string(b)

	return schema, err
}
