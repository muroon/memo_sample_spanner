package usecase

import (
	"memo_sample_spanner/adapter/spanner"
	"memo_sample_spanner/domain/repository"
	"memo_sample_spanner/domain/transaction"
	"memo_sample_spanner/infra/error"
	"memo_sample_spanner/testutil"
)

var testManager testutil.TestManager

func init() {
	testManager = testutil.NewTestManager()
}

// getSpannerRepository spanner repository
func getSpannerRepository() (
	repository.MemoRepository, repository.TagRepository, apperror.ErrorManager,
) {
	return testManager.GetSpannerRepository()
}

// getSpannerTransaction get transaction
func getSpannerTransaction() transaction.ITransaction {
	return spanner.NewTransaction()
}

// connectTestDB DB接続
func connectTestDB() {
	if err := testManager.ConnectTestDB(); err != nil {
		panic(err)
	}
}

// closeTestDB DB切断
func closeTestDB() {
	_ = testManager.CloseTestDB()
}
