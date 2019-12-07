package spanner

import (
	"context"
	"fmt"
	"memo_sample_spanner/domain/model"
	"strings"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
)

func TestTagSaveInDBSuccess(t *testing.T) {

	repo := NewTagRepository()

	ctx := context.Background()
	connectTestDB(ctx)
	defer closeTestDB()

	me, err := repo.Save(ctx, "Tag First")
	if err != nil {
		t.Error("failed TestTagSaveInDBSuccess Save", err)
	}

	tag, err := repo.Get(ctx, me.TagID)
	if err != nil {
		t.Error("failed TestTagSaveInDBSuccess Get", err)
	}

	t.Log(tag)

}

func TestTagAndMemoGetAllByMemoIDSuccess(t *testing.T) {

	tx := NewTransaction()
	repoT := NewTagRepository()
	repoM := NewMemoRepository()

	ctx := context.Background()

	connectTestDB(ctx)
	defer closeTestDB()

	var memo *model.Memo
	var tag *model.Tag
	_, err := tx.ReadWriteTransaction(ctx, func(ctx context.Context) (err error) {
		memo, err = repoM.Save(ctx, "GetAllByMemoID Test Memo")
		if err != nil {
			return err
		}

		tag, err = repoT.Save(ctx, "GetAllByMemoID Test Tag")
		if err != nil {
			return err
		}

		err = repoT.SaveTagAndMemo(ctx, tag.TagID, memo.MemoID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("TestTagAndMemoGetAllByMemoIDSuccess targetMemoID:%s", memo.MemoID)

	flag := false
	list, err := repoT.GetAllByMemoID(ctx, memo.MemoID)
	if err != nil {
		t.Error(err)
		return
	}

	for _, tg := range list {
		if tg.TagID == tag.TagID {
			flag = true
			t.Log(tg)
		}
	}

	if !flag {
		t.Error(fmt.Errorf("GetAllByMemoID Error"))
	}
}

func TestTagAndMemoSearchMemoIDsByTitleSuccess(t *testing.T) {

	repoT := NewTagRepository()
	repoM := NewMemoRepository()

	ctx := context.Background()

	connectTestDB(ctx)
	defer closeTestDB()

	memo, err := repoM.Save(ctx, "SearchMemoIDsByTitle Test Memo")
	if err != nil {
		panic(err)
	}

	tag, err := repoT.Save(ctx, "SearchMemoIDsByTitle Test Tag")
	if err != nil {
		panic(err)
	}

	err = repoT.SaveTagAndMemo(ctx, tag.TagID, memo.MemoID)
	if err != nil {
		panic(err)
	}

	list, err := repoT.SearchMemoIDsByTitle(ctx, tag.Title.StringVal)
	if err != nil {
		panic(err)
	}

	flag := false
	for _, id := range list {
		if id == memo.MemoID {
			flag = true
		}
	}
	if !flag {
		t.Error(fmt.Errorf("SearchMemoIDsByTitle Error"))
	}

	memos, tags, err := repoT.SearchMemoAndTagByTagTitle(ctx, tag.Title.StringVal)
	if err != nil {
		panic(err)
	}

	for _, tg := range tags {
		if !strings.Contains(tg.Title.StringVal, tag.Title.StringVal) {
			t.Error(fmt.Errorf("SearchMemoIDsByTitle Error"))
			break
		}
	}

	flag = false
	for _, me := range memos {
		if me.MemoID == memo.MemoID {
			flag = true
		}
	}
	if !flag {
		t.Error(fmt.Errorf("SearchMemoIDsByTitle Error"))
	}
}
