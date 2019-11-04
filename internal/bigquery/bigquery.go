package bigquery

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"time"
)

type TweetEntry struct {
	ID            string
	InReplyToUser string
	RawSource     string
	FullText      string
	Tokens        []string
	Timestamp     time.Time
	Source        string
}

func Schema() bigquery.Schema {
	schema, err := bigquery.InferSchema(TweetEntry{})
	if err != nil {
		panic(fmt.Errorf("cannot infer schema: %w", err))
	}

	return schema
}
