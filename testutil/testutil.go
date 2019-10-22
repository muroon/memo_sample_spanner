package testutil

import (
	"context"
	"memo_sample_spanner/adapter/spanner"
	"memo_sample_spanner/domain/repository"
	"memo_sample_spanner/infra/cloudspanner"
	"memo_sample_spanner/infra/error"
)

// NewTestManager test util
func NewTestManager() TestManager {
	return testManager{}
}

// TestManager test manager
type TestManager interface {
	GetSpannerRepository() (repository.MemoRepository, repository.TagRepository, apperror.ErrorManager)
	ConnectTestDB() error
	CloseTestDB() error
}

// testManager test manager
type testManager struct {
}

func (t testManager) GetSpannerRepository() (
	repository.MemoRepository, repository.TagRepository, apperror.ErrorManager,
) {
	return spanner.NewMemoRepository(), spanner.NewTagRepository(), apperror.NewErrorManager()
}

// connectTestDB DB接続
func (t testManager) ConnectTestDB() error {
	return cloudspanner.OpenClient(context.Background(), false)
}

// closeTestDB DB切断
func (t testManager) CloseTestDB() error {
	cloudspanner.CloseClient()
	return nil
}
