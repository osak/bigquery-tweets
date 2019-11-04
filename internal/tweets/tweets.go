package tweets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type rawTweet struct {
	ID            string      `json:"id"`
	InReplyToUser string      `json:"in_reply_to_screen_name"`
	RawSource     string      `json:"source"`
	FullText      string      `json:"full_text"`
	Timestamp     TwitterTime `json:"created_at"`
}

type Tweet struct {
	rawTweet
	Source string
}

func (t *Tweet) UnmarshalJSON(body []byte) error {
	if err := json.Unmarshal(body, &t.rawTweet); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	// Rip off tag
	start := strings.Index(t.rawTweet.RawSource, ">") + 1
	end := strings.LastIndex(t.rawTweet.RawSource, "<")
	t.Source = t.rawTweet.RawSource[start:end]
	return nil
}

func Load(path string) ([]Tweet, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s': %w", path, err)
	}
	defer f.Close()

	if err := skipHeader(f); err != nil {
		return nil, fmt.Errorf("failed to skip header: %w", err)
	}

	dec := json.NewDecoder(f)
	result := make([]Tweet, 0)
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to read tweets: %w", err)
	}
	return result, nil
}

const TweetJsHeader = "window.YTD.tweet.part0 = "

// skipHeader skips "window.YTD.tweet.part0 = " at the beginning of tweet.js which makes it a invalid JSON.
func skipHeader(f *os.File) error {
	buf := make([]byte, len(TweetJsHeader))
	n, err := f.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read from file: %w", err)
	}
	if n != len(buf) || string(buf) != TweetJsHeader {
		return fmt.Errorf("file does not start with expected header. read='%s'", string(buf))
	}
	return nil
}
