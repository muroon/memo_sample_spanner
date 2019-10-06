//+build wireinject

package di

import (
	"memo_sample_spanner/adapter/db"
	"memo_sample_spanner/adapter/error"
	"memo_sample_spanner/adapter/logger"
	"memo_sample_spanner/adapter/memory"
	view "memo_sample_spanner/adapter/view/render"
	"memo_sample_spanner/interface/api"
	"memo_sample_spanner/usecase"

	"github.com/google/wire"
)

// ProvideAPI inject api using wire
var ProvideAPI = wire.NewSet(
	ProvideUsecaseIterator,
	api.NewAPI,
)

// ProvidePresenter inject presenter using wire
var ProvidePresenter = wire.NewSet(
	ProvideRender,
	ProvideLog,
	api.NewPresenter,
	ProvideErrorManager,
)

// ProvideMemoUsecase inject memo usecase using wire
var ProvideMemoUsecase = wire.NewSet(
	ProvideDBRepository, // or ProvideInMemoryRepository
	usecase.NewMemo,
)

// ProvideUsecaseIterator inject usecase itetator using wire
var ProvideUsecaseIterator = wire.NewSet(
	ProvidePresenter,
	ProvideMemoUsecase,
	usecase.NewInteractor,
)

// ProvideInMemoryRepository inject repository using wire
var ProvideInMemoryRepository = wire.NewSet(
	memory.NewTransactionRepository,
	memory.NewMemoRepository,
	memory.NewTagRepository,
)

// ProvideDBRepository inject repository using wire
var ProvideDBRepository = wire.NewSet(
	db.NewTransactionRepository,
	db.NewMemoRepository,
	db.NewTagRepository,
)

// ProvideLog inject log using wire
var ProvideLog = wire.NewSet(loggersub.NewLogger)

// ProvideRender inject render using wire
var ProvideRender = wire.NewSet(view.NewJSONRender)

// ProvideErrorManager inject error manager using wire
var ProvideErrorManager = wire.NewSet(apperrorsub.NewErrorManager)

// InjectAPIServer build inject api using wire
func InjectAPIServer() api.API {
	wire.Build(
		ProvideAPI,
	)
	return nil
}
