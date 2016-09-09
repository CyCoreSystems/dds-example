package dds

// NewModel creates a new model factory
func NewModel(name string, f func() interface{}) *ModelFactory {
	return &ModelFactory{name, f}
}

// A ModelFactory is a builder for a distributed data model
type ModelFactory struct {
	EntityName string
	Build      func() interface{}
}
