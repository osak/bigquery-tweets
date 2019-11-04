package tweets

import (
	"encoding/json"
	"fmt"
	"time"
)

type TwitterTime struct {
	time.Time
}

func (tt *TwitterTime) UnmarshalJSON(body []byte) error {
	var timestamp string
	if err := json.Unmarshal(body, &timestamp); err != nil {
		return fmt.Errorf("cannot unmarshal timestamp '%s': %w", body, err)
	}

	t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", timestamp)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp '%s': %w", timestamp, err)
	}
	tt.Time = t
	return nil
}

