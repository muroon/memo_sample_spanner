package db

import (
	"context"
	"memo_sample_spanner/domain/repository"
	"testing"
)

func getMemoRepositoryForTest() repository.MemoRepository {
	return NewMemoRepository()
}

func TestMemoSaveInDBSuccess(t *testing.T) {

	repo := getMemoRepositoryForTest()

	ctx := context.Background()

	connectTestDB()
	defer closeTestDB()

	// 1件名
	memo, err := repo.Save(ctx, "First")
	if err != nil {
		t.Error("failed TestMemoSaveInMemorySuccess Save", err)
	}

	memoGet, err := repo.Get(ctx, memo.ID)
	if err != nil || memoGet.ID != memo.ID {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, memoGet.ID)
	}

	// 2件名
	memo, err = repo.Save(ctx, "Second")
	if err != nil {
		t.Error("failed TestMemoSaveInMemorySuccess GenerateID", err)
	}

	memoGet, err = repo.Get(ctx, memo.ID)
	if err != nil || memoGet.ID != memo.ID {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, memoGet.ID)
	}

	//　全件取得
	list, err := repo.GetAll(ctx)
	if err != nil {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, len(list))
	}

	for _, v := range list {
		t.Logf("TestMemoSaveInMemorySuccess GetAll MemoRepository id:%d, text:%s", v.ID, v.Text)
	}
}

func TestMemoTransactionCommitSuccess(t *testing.T) {

	repo := getMemoRepositoryForTest()
	repoTx := getTransactionRepositoryForTest()

	ctx := context.Background()

	connectTestDB()
	defer func() {
		if err := recover(); err != nil {
			repoTx.Rollback(ctx)
			t.Error(err)
		}
		closeTestDB()
	}()

	ctx, err := repoTx.Begin(ctx)
	if err != nil {
		panic(err)
	}

	_, err = repo.Save(ctx, "Transaction Commit Test")
	if err != nil {
		panic(err)
	}

	_, err = repoTx.Commit(ctx)
	if err != nil {
		panic(err)
	}
}

func TestMemoTransactionRollbackSuccess(t *testing.T) {

	repo := getMemoRepositoryForTest()
	repoTx := getTransactionRepositoryForTest()

	ctx := context.Background()

	connectTestDB()
	defer func() {
		if err := recover(); err != nil {
			repoTx.Rollback(ctx)
			t.Error(err)
		}
		closeTestDB()
	}()

	ctx, err := repoTx.Begin(ctx)
	if err != nil {
		panic(err)
	}

	_, err = repo.Save(ctx, "Transaction Rollback Test")
	if err != nil {
		panic(err)
	}

	repoTx.Rollback(ctx)
}

func TestMemoSearchSuccess(t *testing.T) {

	repo := getMemoRepositoryForTest()

	ctx := context.Background()

	connectTestDB()
	defer closeTestDB()

	word := "Memo Search Test"
	_, err := repo.Save(ctx, word)
	if err != nil {
		t.Error(err)
	}

	word = "Memo"
	list, err := repo.Search(ctx, word)
	if err != nil {
		t.Error(err)
	}

	for _, m := range list {
		t.Log(m)
	}
}

func TestMemoGetAllByIDsSuccess(t *testing.T) {

	repo := getMemoRepositoryForTest()

	ctx := context.Background()

	connectTestDB()
	defer closeTestDB()

	word := "Dummy First"
	memo1, err := repo.Save(ctx, word)
	if err != nil {
		t.Error(err)
	}

	word = "Dummy Second"
	memo2, err := repo.Save(ctx, word)
	if err != nil {
		t.Error(err)
	}

	ids := []int{memo1.ID, memo2.ID}
	list, err := repo.GetAllByIDs(ctx, ids)
	if err != nil {
		t.Error(err)
	}

	for _, m := range list {
		t.Log(m)
	}
}
