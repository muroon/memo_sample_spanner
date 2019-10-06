package main

import (
	"memo_sample_spanner/di"
	"memo_sample_spanner/infra/database"
	"net/http"
)

func main() {
	(*database.GetDBM()).ConnectDB()
	defer (*database.GetDBM()).CloseDB()

	api := di.InjectAPIServer()
	http.HandleFunc("/", api.GetMemos)
	http.HandleFunc("/post", api.PostMemo)
	http.HandleFunc("/post/memo_tags", api.PostMemoAndTags)
	http.HandleFunc("/search/tags_memos", api.SearchTagsAndMemos)
	http.ListenAndServe(":8080", nil)
}
