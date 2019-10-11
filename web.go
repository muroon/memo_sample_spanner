package main

import (
	"context"
	"memo_sample_spanner/di"
	"memo_sample_spanner/infra/cloudspanner"
	"net/http"
)

func main() {
	cloudspanner.OpenClient(context.Background())
	defer cloudspanner.CloseClient()

	api := di.InjectAPIServer()
	http.HandleFunc("/", api.GetMemos)
	http.HandleFunc("/post", api.PostMemo)
	http.HandleFunc("/post/memo_tags", api.PostMemoAndTags)
	http.HandleFunc("/search/tags_memos", api.SearchTagsAndMemos)
	http.ListenAndServe(":8080", nil)
}
