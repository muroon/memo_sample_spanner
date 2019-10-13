package render

import (
	"fmt"
	"memo_sample_spanner/domain/model"
	"memo_sample_spanner/view/model/json"
)

// NewJSONRender new json render
func NewJSONRender() JSONRender {
	return jsonRender{}
}

type jsonRender struct {
}

func (m jsonRender) ConvertMemo(md *model.Memo) *json.Memo {
	mj := &json.Memo{
		MemoID: md.MemoID,
		Text:   md.Text.StringVal,
	}
	return mj
}

func (m jsonRender) ConvertMemos(list []*model.Memo) []*json.Memo {
	listJSON := []*json.Memo{}
	for _, v := range list {
		listJSON = append(listJSON, m.ConvertMemo(v))
	}
	return listJSON
}

func (m jsonRender) ConvertTag(md *model.Tag) *json.Tag {
	mj := &json.Tag{
		TagID: md.TagID,
		Title: md.Title.StringVal,
	}
	return mj
}

func (m jsonRender) ConvertTags(list []*model.Tag) []*json.Tag {
	listJSON := []*json.Tag{}
	for _, v := range list {
		listJSON = append(listJSON, m.ConvertTag(v))
	}
	return listJSON
}

func (m jsonRender) ConvertPostMemoAndTagsResult(memo *model.Memo, tags []*model.Tag) *json.PostMemoAndTagsResult {

	return &json.PostMemoAndTagsResult{
		Memo: m.ConvertMemo(memo),
		Tags: m.ConvertTags(tags),
	}
}

func (m jsonRender) ConvertSearchTagsAndMemosResult(memos []*model.Memo, tags []*model.Tag) *json.SearchTagsAndMemosResult {

	return &json.SearchTagsAndMemosResult{
		Tags:  m.ConvertTags(tags),
		Memos: m.ConvertMemos(memos),
	}
}

func (m jsonRender) ConvertError(err error, code int) *json.Error {
	mess := fmt.Sprintf("API: %T(%v)\n", err, err)

	return &json.Error{Code: code, Msg: mess}
}
