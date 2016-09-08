package dnats

import (
	"dds"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

// ListenHandler listens for an action handler
func ListenHandler(ctx context.Context, af *dds.ActionFactory, handler dds.Handler) error {

	actionName := af.ActionName

	var nc *nats.Conn
	var err error

	for i := 0; i != 3 && nc == nil; i++ {
		<-time.After(500 * time.Millisecond)
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			fmt.Printf("Error connecting to nats: '%v'", err)
		}
	}

	if nc == nil {
		return errors.New("Failed to connect to nats, giving up\n")
	}

	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	go func() {
		<-ctx.Done()
		nc.Close()
	}()

	c.Subscribe(actionName, func(subj, reply string, data []byte) {
		req := af.RequestBuild()

		if err := json.Unmarshal(data, &req); err != nil {
			panic(err)
		}
		resp, _ := handler.Handle(req)
		c.Publish(reply, resp)
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

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	return func(i interface{}) (resp interface{}, err error) {
		replyID := uuid.NewV1().String()

		ch := make(chan string)
		defer close(ch)

		c.Subscribe(replyID, func(body []byte) {
			resp = af.ResponseBuild()
			err = json.Unmarshal(body, resp)
			ch <- ""
		})

		if err := c.PublishRequest(actionName, replyID, i); err != nil {
			return nil, err
		}

		select {
		case <-ch:
			return
		case <-time.After(2 * time.Second):
			return nil, errors.New("timeout error")
		}

	}
}
