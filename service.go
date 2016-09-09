package dds

import "golang.org/x/net/context"

// A DataService is a process that listens for data operation requests
type DataService struct {
	ctx    context.Context
	cancel context.CancelFunc

	storage Storage

	transport Transport

	modelFactory *ModelFactory
}

// NewDataService creates a new service
func NewDataService(mf *ModelFactory, storage Storage, transport Transport) *DataService {
	svc := &DataService{}

	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	svc.storage = storage
	svc.transport = transport
	svc.modelFactory = mf

	return svc
}

// Listen listens for requests
func (svc *DataService) Listen() error {
	return svc.transport.Model(svc.modelFactory, svc.storage)
}

// Stop stops the service
func (svc *DataService) Stop() {
	svc.transport.Close()
}
