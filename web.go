package main

import (
	"context"
	"flag"
	"memo_sample_spanner/di"
	"memo_sample_spanner/infra/cloudspanner"
	"memo_sample_spanner/infra/logger"
	"net/http"
)

func main() {
	local := flag.Bool("local", false, "use local spanner")
	flag.Parse()
	logger.NewLogger().Infof("Start local:%v\n", *local)

	err := cloudspanner.OpenClient(context.Background(), *local)
	defer cloudspanner.CloseClient()
	if err != nil {
		logger.NewLogger().Errorf("db open error: %#+v\n", err)
	}

	api := di.InjectAPIServer()
	http.HandleFunc("/", api.GetMemos)
	http.HandleFunc("/post", api.PostMemo)
	http.HandleFunc("/post/memo_tags", api.PostMemoAndTags)
	http.HandleFunc("/search/tags_memos", api.SearchTagsAndMemos)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.NewLogger().Errorf("ListenAndServe error: %#+v\n", err)
	}
}
