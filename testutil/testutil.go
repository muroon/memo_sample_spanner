package testutil

import (
	"memo_sample_spanner/adapter/db"
	"memo_sample_spanner/adapter/error"
	"memo_sample_spanner/adapter/memory"
	"memo_sample_spanner/domain/repository"
	"memo_sample_spanner/infra/database"
	"memo_sample_spanner/infra/error"
)

// NewTestManager test util
func NewTestManager() TestManager {
	return testManager{}
}

// TestManager test manager
type TestManager interface {
	GgetInMemoryRepository() (repository.TransactionRepository, repository.MemoRepository, repository.TagRepository, apperror.ErrorManager)
	GetDBRepository() (repository.TransactionRepository, repository.MemoRepository, repository.TagRepository, apperror.ErrorManager)
	ConnectTestDB() error
	CloseTestDB() error
}

// testManager test manager
type testManager struct {
}

func (t testManager) GgetInMemoryRepository() (repository.TransactionRepository, repository.MemoRepository, repository.TagRepository, apperror.ErrorManager) {
	return memory.NewTransactionRepository(), memory.NewMemoRepository(), memory.NewTagRepository(), apperrorsub.NewErrorManager()
}

func (t testManager) GetDBRepository() (repository.TransactionRepository, repository.MemoRepository, repository.TagRepository, apperror.ErrorManager) {
	return db.NewTransactionRepository(), db.NewMemoRepository(), db.NewTagRepository(), apperrorsub.NewErrorManager()
}

// connectTestDB DB接続
func (t testManager) ConnectTestDB() error {
	return (*database.GetDBM()).ConnectTestDB()
}

// closeTestDB DB切断
func (t testManager) CloseTestDB() error {
	return (*database.GetDBM()).CloseDB()
}
