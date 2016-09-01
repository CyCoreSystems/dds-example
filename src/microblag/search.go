package microblag

import "time"

// A Search is a request to search for a series of entries
type Search struct {
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	User    *string   `json:"user_id"`
}
