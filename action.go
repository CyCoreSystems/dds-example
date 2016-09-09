package dds

import "golang.org/x/net/context"

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

// An ActionService is a process that listens for RPC operations
type ActionService struct {
	ctx    context.Context
	cancel context.CancelFunc

	transport Transport

	handler       Handler
	actionFactory *ActionFactory
}

// NewActionService creates a new service
func NewActionService(af *ActionFactory, handler Handler, transport Transport) *ActionService {
	svc := &ActionService{}

	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	svc.transport = transport
	svc.actionFactory = af
	svc.handler = handler

	return svc
}

// Listen listens for requests
func (svc *ActionService) Listen() error {
	return svc.transport.Action(svc.actionFactory, svc.handler)
}

// Stop stops the service
func (svc *ActionService) Stop() {
	svc.transport.Close()
}
