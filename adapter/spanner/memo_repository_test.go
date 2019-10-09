package spanner

import (
	"context"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
)

func TestMemoSaveInMemorySuccess(t *testing.T) {
	ctx := context.Background()

	connectTestDB(ctx)
	defer closeTestDB()

	repo := NewMemoRepository()

	// 1件名
	memo, err := repo.Save(ctx, "First")
	if err != nil {
		t.Error("failed TestMemoSaveInMemorySuccess Save", err)
	}

	memoGet, err := repo.Get(ctx, memo.MemoID)
	if err != nil || memoGet.MemoID != memo.MemoID {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, memoGet.MemoID)
	}
	t.Logf("TestMemoSaveInMemorySuccess Get MemoRepository id:%s, text:%s", memoGet.MemoID, memoGet.Text)

	// 2件名
	memo, err = repo.Save(ctx, "Second")
	if err != nil {
		t.Error("failed TestMemoSaveInMemorySuccess Save", err)
	}

	memoGet, err = repo.Get(ctx, memo.MemoID)
	if err != nil || memoGet.MemoID != memo.MemoID {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, memoGet.MemoID)
	}
	t.Logf("TestMemoSaveInMemorySuccess Get MemoRepository id:%s, text:%s", memoGet.MemoID, memoGet.Text)

	//　全件取得
	list, err := repo.GetAll(ctx)
	if err != nil || len(list) < 2 {
		t.Error("failed TestMemoSaveInMemorySuccess Get", err, len(list))
	}

	for _, v := range list {
		t.Logf("TestMemoSaveInMemorySuccess GetAll MemoRepository id:%s, text:%s", v.MemoID, v.Text)
	}
}

func TestMemoSearchSuccess(t *testing.T) {

	ctx := context.Background()

	connectTestDB(ctx)
	defer closeTestDB()

	repo := NewMemoRepository()

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

	repo := NewMemoRepository()

	ctx := context.Background()

	connectTestDB(ctx)
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

	ids := []string{memo1.MemoID, memo2.MemoID}
	list, err := repo.GetAllByIDs(ctx, ids)
	if err != nil {
		t.Error(err)
	}

	for _, m := range list {
		t.Log(m)
	}
}
