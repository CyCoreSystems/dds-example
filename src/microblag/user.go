package microblag

import (
	"dds"
	"fmt"
	"time"
)

// UserFactory is a factory for distributed data model binding
var UserFactory = dds.NewModel("microblag-user", func() interface{} { return &User{} })

// A User is a user of a system
type User struct {
	ID          string    `json:"id,omitempty"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

func (u *User) String() string {

	dp := "<none>"
	if u.DisplayName != nil {
		dp = *u.DisplayName
	}

	return fmt.Sprintf("User{%s %s %v %v %v}", u.ID, u.Username, dp, u.Created, u.Updated)
}
