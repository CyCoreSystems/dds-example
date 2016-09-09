package microblag

import (
	"time"

	"github.com/CyCoreSystems/dds"
)

// EntryFactory is a factory for distributed data model binding
var EntryFactory = dds.NewModel("microblag-entry", func() interface{} { return &Entry{} })

// An Entry is a string published by a user
type Entry struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	Text    string    `json:"text"`
	Created time.Time `json:"created"`

	Parent *string `json:"parent_id"`
}
