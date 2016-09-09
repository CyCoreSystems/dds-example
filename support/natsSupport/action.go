package natsSupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
)

// ListenHandler listens for an action handler
func (nt *natsTransport) Action(af *dds.ActionFactory, handler dds.Handler) error {

	actionName := af.ActionName

	nc := nt.nc

	var reporter = func(err error, msg string) {
		fmt.Printf("ERR: %s: '%v'", msg, err)
	}

	subscribe(nc, reporter, actionName, func(subj string, data []byte, reply Reply) {
		req := af.RequestBuild()
		if err := json.Unmarshal(data, &req); err != nil {
			reply(nil, err)
			return
		}

		resp, err := handler.Handle(req)
		reply(resp, err)
	})

	return nil

}

// ActionClient builds a client that can invoke a remote action
func ActionClient(af *dds.ActionFactory) dds.Action {
	actionName := af.ActionName

	var nc *nats.Conn
	var err error

	for i := 0; i != 3 && nc == nil; i++ {
		<-time.After(500 * time.Millisecond)
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			fmt.Printf("Error connecting to nats: '%v'\n", err)
		}
	}

	if nc == nil {
		panic(errors.New("Failed to connect to nats, giving up\n"))
	}

	return func(i interface{}) (resp interface{}, err error) {
		err = request(nc, actionName, i, &resp)
		return
	}
}
