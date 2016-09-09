package dds

/*

	Usage:
		nt, _ := nats.NewTransport()
		ctx := NewContext(nt)
		svc := ctx.RPC.Service(MyRPCMethod, MyHandler)
		svc.Listen()


		client := ctx.RPC.Client(MyRPCMethod)
		resp, err := client(request)

*/

// A Context is a wrapper around a transort, used to build clients and servers
type Context struct {
	Data Data
	RPC  RPC

	transport Transport
}

func NewContext(transport Transport) *Context {
	return &Context{
		transport: transport,
		Data:      &transportData{transport},
		RPC:       &transportRPC{transport},
	}
}

type Data interface {
	Client(*ModelFactory) Model
	Service(*ModelFactory, Storage) *DataService
}

type RPC interface {
	Client(*ActionFactory) Action
	Service(*ActionFactory, Handler) *ActionService
}

type transportData struct {
	transport Transport
}

func (td *transportData) Client(mf *ModelFactory) Model {
	panic("Not implemented")
}

func (td *transportData) Service(mf *ModelFactory, st Storage) *DataService {
	return NewDataService(mf, st, td.transport)
}

type transportRPC struct {
	transport Transport
}

func (rpc *transportRPC) Client(af *ActionFactory) Action {
	panic("Not implemented")
}

func (rpc *transportRPC) Service(af *ActionFactory, handler Handler) *ActionService {
	return NewActionService(af, handler, rpc.transport)
}
