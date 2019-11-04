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

type DumpableSchemaEntry struct {
	Name string             `json:"name"`
	Mode string             `json:"mode"`
	Type bigquery.FieldType `json:"type"`
}
type DumpableSchema []DumpableSchemaEntry

func Schema() DumpableSchema {
	schema, err := bigquery.InferSchema(TweetEntry{})
	if err != nil {
		panic(fmt.Errorf("cannot infer schema: %w", err))
	}

	result := make(DumpableSchema, len(schema))
	for i, s := range schema {
		var mode string
		if s.Repeated {
			mode = "REPEATED"
		} else if s.Required {
			mode = "REQUIRED"
		} else {
			mode = "NULLABLE"
		}
		result[i] = DumpableSchemaEntry{
			Name: s.Name,
			Mode: mode,
			Type: s.Type,
		}
	}
	return result
}
