package dds

// NewAction builds a new action factory
func NewAction(name string, req func() interface{}, resp func() interface{}) *ActionFactory {
	return &ActionFactory{name, req, resp}
}

// An ActionFactory is a builder for a distributed action
type ActionFactory struct {
	ActionName    string
	RequestBuild  func() interface{}
	ResponseBuild func() interface{}
}

// An Action is an RPC style operation
type Action func(i interface{}) (interface{}, error)
