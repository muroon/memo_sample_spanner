package json

// Memo Memo's Out Entity
type Memo struct {
	MemoID string `json:"memo_id"`
	Text   string `json:"text"`
}

// Tag Tag's Out Entity
type Tag struct {
	TagID string `json:"tag_id"`
	Title string `json:"title"`
}

// PostMemoAndTagsResult Out Entity For PostMemoAndTags Result
type PostMemoAndTagsResult struct {
	Memo *Memo  `json:"memo"`
	Tags []*Tag `json:"tags"`
}

// SearchTagsAndMemosResult Out Entity For SearchTagsAndMemos Result
type SearchTagsAndMemosResult struct {
	Tags  []*Tag  `json:"tags"`
	Memos []*Memo `json:"memos"`
}
