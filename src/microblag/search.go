package microblag

import (
	"dds"
	"time"
)

// Search is the distributed action for searching
var Search = dds.NewAction("search", func() interface{} { return &SearchRequest{} }, func() interface{} { return &[]string{} })

// A SearchRequest is a request to search for a series of entries
type SearchRequest struct {
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	User    *string   `json:"user_id"`
}
