package microblag

import (
	"dds"
	"time"
)

// UserFactory is a factory for distributed data model binding
var UserFactory = dds.NewModel(func() interface{} { return &User{} })

// A User is a user of a system
type User struct {
	ID          string    `json:"id,omitempty"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
