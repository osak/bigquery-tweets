package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/osak/bigquery-tweets/internal/bigquery"
	"github.com/osak/bigquery-tweets/internal/tokenizer"
	"github.com/osak/bigquery-tweets/internal/tweets"
	"os"
)

type params struct {
	inputPath  string
	outputPath string
	schemaPath string
}

func main() {
	params := params{}
	flag.StringVar(&params.inputPath, "input", "", "Path to tweets.csv")
	flag.StringVar(&params.outputPath, "output", "", "Path to save parsed csv")
	flag.StringVar(&params.schemaPath, "schema", "", "Path to save schema file")
	flag.Parse()

	if params.inputPath == "" || params.outputPath == "" || params.schemaPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	allTweets, err := tweets.Load(params.inputPath)
	if err != nil {
		panic(err)
	}
	t := tokenizer.New()

	entries := make([]bigquery.TweetEntry, len(allTweets))
	for i, tweet := range allTweets {
		tokens := t.Tokenize(tweet.FullText)
		entries[i] = bigquery.TweetEntry{
			ID:            tweet.ID,
			InReplyToUser: tweet.InReplyToUser,
			Source:        tweet.Source,
			FullText:      tweet.FullText,
			Tokens:        tokens,
			Timestamp:     tweet.Timestamp.Time,
		}
	}
	if err = saveEntries(params.outputPath, entries); err != nil {
		panic(err)
	}
	if err = saveSchema(params.schemaPath); err != nil {
		panic(err)
	}
}

// saveEntries saves given entries as newline-delimited JSON.
func saveEntries(path string, entries []bigquery.TweetEntry) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot open file '%s': %w", path, err)
	}

	w := json.NewEncoder(out)
	for _, entry := range entries {
		if err := w.Encode(entry); err != nil {
			return fmt.Errorf("cannot write entry (ID=%s): %w", entry.ID, err)
		}
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("cannot properly close file: %w", err)
	}
	return nil
}

func saveSchema(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot open file '%s': %w", path, err)
	}

	w := json.NewEncoder(out)
	schema := bigquery.Schema()
	err = w.Encode(schema)
	if err != nil {
		return fmt.Errorf("cannot write to file '%s': %w", path, err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("cannot properly close file: %w", err)
	}
	return nil
}
