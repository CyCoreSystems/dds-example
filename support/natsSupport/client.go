package natsSupport

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	uuid "src/github.com/satori/go.uuid"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
)

// Client creates a client
func Client(mf *dds.ModelFactory) dds.Model {

	entityName := mf.EntityName

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

	RegisterEncoder(mf, entityName)

	c, err := nats.NewEncodedConn(nc, "dds."+entityName)
	if err != nil {
		panic(err)
	}

	return &natsModel{
		c:          c,
		entityName: entityName,
	}
}

type natsModel struct {
	c          *nats.EncodedConn
	entityName string
}

func (m *natsModel) Close() {
	m.c.Close()
}

func (m *natsModel) Get(id string, i interface{}) error {
	replyID := uuid.NewV1().String()

	ch := make(chan error)
	defer close(ch)

	m.c.Subscribe(replyID, func(resp *container) {
		if resp.I == nil {
			ch <- resp.Err
			return
		}
		reflect.ValueOf(i).Elem().Set(reflect.ValueOf(resp.I))
		ch <- nil
	})

	if err := m.c.PublishRequest(m.entityName+".get."+id, replyID, i); err != nil {
		return err
	}

	select {
	case err := <-ch:
		return err
	case <-time.After(2 * time.Second):
		return errors.New("timeout error")
	}
}

func (m *natsModel) Create(i interface{}) (string, error) {
	replyID := uuid.NewV1().String()

	ch := make(chan string)
	defer close(ch)

	m.c.Subscribe(replyID, func(resp map[string]interface{}) {
		id := resp["id"].(string)
		ch <- id
	})

	if err := m.c.PublishRequest(m.entityName+".create", replyID, i); err != nil {
		return "", err
	}

	select {
	case id := <-ch:
		return id, nil
	case <-time.After(2 * time.Second):
		return "", errors.New("timeout error")
	}
}

func (m *natsModel) Delete(ID string) error {
	replyID := uuid.NewV1().String()

	ch := make(chan string)
	defer close(ch)

	m.c.Subscribe(replyID, func(_, _ string, d map[string]interface{}) {
		ch <- d["status"].(string)
	})

	if err := m.c.PublishRequest(m.entityName+".delete."+ID, replyID, ""); err != nil {
		return err
	}

	select {
	case <-ch: //TODO: check resp
		return nil
	case <-time.After(2 * time.Second):
		return errors.New("timeout error")
	}
}

func (m *natsModel) Update(ID string, i interface{}) error {
	replyID := uuid.NewV1().String()

	ch := make(chan string)
	defer close(ch)

	m.c.Subscribe(replyID, func(_, _ string, d map[string]interface{}) {
		ch <- d["status"].(string)
	})

	if err := m.c.PublishRequest(m.entityName+".update", replyID, i); err != nil {
		return err
	}

	select {
	case <-ch: //TODO: check resp
		return nil
	case <-time.After(2 * time.Second):
		return errors.New("timeout error")
	}
}

func (m *natsModel) Subscribe(evt string) dds.Subscription {
	var nc natsSubscription
	nc.c = make(chan dds.Event, 10)
	sx, _ := m.c.Subscribe("events."+m.entityName+"."+evt, func(subj string, _ string, data []byte) {

		items := strings.Split(subj, ".")

		nc.c <- dds.Event{
			Entity: m.entityName,
			Type:   items[2],
			Metadata: map[string]interface{}{
				"data": string(data),
			},
		}
	})

	nc.sx = sx

	return &nc
}

type natsSubscription struct {
	c  chan dds.Event
	sx *nats.Subscription
}

func (nc *natsSubscription) C() <-chan dds.Event {
	return nc.c
}

func (nc *natsSubscription) Close() {
	nc.sx.Unsubscribe()
	if nc.c != nil {
		close(nc.c)
	}
	nc.c = nil
}
