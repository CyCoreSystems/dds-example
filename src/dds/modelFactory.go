package dds

// NewModel creates a new model factory
func NewModel(f func() interface{}) *ModelFactory {
	return &ModelFactory{f}
}

// A ModelFactory is a builder for a distributed data model
type ModelFactory struct {
	Build func() interface{}
}

// Service runs the distributed data model service given the
// backend and transport
func (mf *ModelFactory) Service(storage Storage) {

}

// Client builds a client that can interact with a distributed data model
func (mf *ModelFactory) Client() Model {
	return nil
}
