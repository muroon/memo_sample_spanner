//+build wireinject

package di

import (
	"memo_sample_spanner/adapter/logger"
	"memo_sample_spanner/adapter/spanner"
	view "memo_sample_spanner/adapter/view/render"
	"memo_sample_spanner/infra/error"
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
	ProvideSpannerRepository,
	usecase.NewMemo,
)

// ProvideUsecaseIterator inject usecase itetator using wire
var ProvideUsecaseIterator = wire.NewSet(
	ProvidePresenter,
	ProviderTransaction,
	ProvideMemoUsecase,
	usecase.NewInteractor,
)

var ProviderTransaction = wire.NewSet(
	spanner.NewTransaction,
)

var ProvideSpannerRepository = wire.NewSet(
	spanner.NewMemoRepository,
	spanner.NewTagRepository,
)

// ProvideLog inject log using wire
var ProvideLog = wire.NewSet(loggersub.NewLogger)

// ProvideRender inject render using wire
var ProvideRender = wire.NewSet(view.NewJSONRender)

// ProvideErrorManager inject error manager using wire
var ProvideErrorManager = wire.NewSet(apperror.NewErrorManager)

// InjectAPIServer build inject api using wire
func InjectAPIServer() api.API {
	wire.Build(
		ProvideAPI,
	)
	return nil
}
