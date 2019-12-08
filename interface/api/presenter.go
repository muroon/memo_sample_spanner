package api

import (
	"context"
	"encoding/json"
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/infra/error"
	"memo_sample_spanner/infra/logger"
	"memo_sample_spanner/usecase"
	"memo_sample_spanner/view/render"
	"net/http"
)

// NewPresenter new presenter
func NewPresenter(render render.JSONRender, log logger.Logger, errm apperror.ErrorManager) usecase.Presenter {
	return presenter{render, log, errm}
}

type presenter struct {
	render render.JSONRender
	log    logger.Logger
	errm   apperror.ErrorManager
}

func (m presenter) ViewMemo(ctx context.Context, md *model.Memo) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertMemo(md))
}

func (m presenter) ViewMemoList(ctx context.Context, list []*model.Memo) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertMemos(list))
}

func (m presenter) ViewTag(ctx context.Context, md *model.Tag) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertTag(md))
}

func (m presenter) ViewTagList(ctx context.Context, list []*model.Tag) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertTags(list))
}

func (m presenter) ViewPostMemoAndTagsResult(ctx context.Context, memo *model.Memo, tags []*model.Tag) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertPostMemoAndTagsResult(memo, tags))
}

func (m presenter) ViewSearchTagsAndMemosResult(ctx context.Context, memos []*model.Memo, tags []*model.Tag) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.JSON(ctx, w, m.render.ConvertSearchTagsAndMemosResult(memos, tags))
}

func (m presenter) ViewError(ctx context.Context, err error) {
	defer deleteResponseWriter(ctx)
	w := getResponseWriter(ctx)

	m.log.Errorf("API: %s\n", m.errm.LogMessage(err))

	m.JSON(ctx, w, m.render.ConvertError(err, m.errm.Code(err)))
}

// JSON render json format
func (m presenter) JSON(ctx context.Context, w http.ResponseWriter, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
