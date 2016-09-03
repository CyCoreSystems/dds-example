package dds

// Events
const (
	AllEvents   = "*"
	CreateEvent = "create"
	UpdateEvent = "update"
	DeleteEvent = "delete"
)

// An Event is a data event
type Event struct {
	Entity   string
	Type     string
	Metadata map[string]interface{}
}
