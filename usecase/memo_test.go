package usecase

import (
	"context"
	"encoding/json"
	"memo_sample_spanner/usecase/input"
	"strings"
	"testing"
)

func TestMemoPostAndGetMemoInDBSuccess(t *testing.T) {
	ctx := context.Background()

	// Spannerでテストした場合
	memo := NewMemo(getSpannerRepository())

	connectTestDB()
	defer closeTestDB()

	text := "Next Memo"

	ipt := &input.PostMemo{Text: text}

	id, err := memo.Post(ctx, *ipt)
	if err != nil {
		t.Error(err)
	}

	iptf := &input.GetMemo{ID: id}
	m, err := memo.GetMemo(ctx, *iptf)
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("TestMemoPostAndGetSuccess Get MemoRepository json: %s", b)

	l, err := memo.GetAllMemoList(ctx)
	if err != nil {
		t.Error(err)
	}
	lb, err := json.Marshal(l)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("TestMemoPostAndGetSuccess GetAllJSON MemoRepository json: %s", lb)
}

func TestMemoSearchTagsAndMemosSuccess(t *testing.T) {
	ctx := context.Background()

	memo := NewMemo(getSpannerRepository())
	connectTestDB()
	defer closeTestDB()

	tx := getSpannerTransaction()

	// test deta post
	memoTexts := []string{"SearchTagsAndMemos 1", "SearchTagsAndMemos 2"}
	tagTitle := "SearchTagsAndMemos"
	tagTitles := []string{tagTitle}
	_, err := tx.ReadWriteTransaction(ctx, func(ctx context.Context) error {
		for _, memoText := range memoTexts {
			ipt1 := &input.PostMemoAndTags{MemoText: memoText, TagTitles: tagTitles}
			_, _, err := memo.PostMemoAndTags(ctx, *ipt1)
			if err != nil {
				t.Error(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	ipt := &input.SearchTagsAndMemos{TagTitle: tagTitle}

	mos, tgs, err := memo.SearchTagsAndMemos(ctx, *ipt)
	if err != nil {
		t.Error(err)
	}

	// check Tag
	for _, tag := range tgs {
		if !strings.Contains(tag.Title.StringVal, tagTitle) {
			t.Errorf("Tag And Memo Save Error. tag.Title:%s", tag.Title)
		}
	}

	// check Memo
	ok := []int{}
	for _, mm := range mos {
		for _, memoText := range memoTexts {
			if mm.Text.StringVal == memoText {
				ok = append(ok, 1)
			}
		}
	}

	if len(ok) < 2 {
		t.Error("Tag And Memo Save Error")
	}
}

func TestReadWriteError(t *testing.T) {
	ctx := context.Background()

	// memo usecase
	memo := NewMemo(getSpannerRepository())

	connectTestDB()
	defer closeTestDB()

	tx := getSpannerTransaction()

	var (
		memoID1 string
		memoID2 string
	)

	_, err := tx.ReadWriteTransaction(ctx, func(ctx context.Context) error {

		var err error
		ipt := input.PostMemo{Text:"NG Form RollBack (Before)."}
		memoID1, err = memo.Post(ctx, ipt)
		if err != nil {
			return err
		}

		err = memo.ValidatePost(*new(input.PostMemo))
		if err != nil {
			t.Log("Validate error")
			return err
		}

		ipt.Text = "NG Form RollBack (After)."
		memoID2, err = memo.Post(ctx, ipt)
		return err
	})
	if err != nil {
		t.Logf("OK. %T, %#+v\n", err, err)
	} else {
		t.Error("ReadWriteError is Error")
	}

	_, err = memo.GetMemo(ctx, input.GetMemo{ID: memoID1})
	if err != nil {
		t.Logf("OK. %T, %#+v\n", err, err)
	} else {
		t.Error("ReadWriteError is Error")
	}

	_, err = memo.GetMemo(ctx, input.GetMemo{ID: memoID2})
	if err != nil {
		t.Logf("OK. %T, %#+v\n", err, err)
	} else {
		t.Error("ReadWriteError is Error")
	}

	// retry
	_, err = tx.ReadWriteTransaction(ctx, func(ctx context.Context) error {
		ipt := input.PostMemo{Text:"Retry Test Memo."}
		_, err = memo.Post(ctx, ipt)
		return err
	})
	if err != nil {
		t.Error(err)
	}
}
