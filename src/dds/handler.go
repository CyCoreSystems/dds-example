package dds

// A Handler is a handler for an RCP action
type Handler interface {
	Handle(i interface{}) (interface{}, error)
}
