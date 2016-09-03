package nats

import (
	"dds"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/nats-io/nats"
	"github.com/nats-io/nats/encoders/builtin"
	"github.com/satori/go.uuid"
)

type statusResult map[string]interface{}

type natsTypeEncoder struct {
	mf *dds.ModelFactory
	je builtin.JsonEncoder
}

func (n *natsTypeEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	switch v.(type) {
	case statusResult:
		b, err := n.je.Encode(subject, v)
		b = append([]byte{2}, b...)
		return b, err
	default:
		return n.je.Encode(subject, v)
	}
}

func (n *natsTypeEncoder) Decode(subject string, data []byte, vPtr interface{}) error {

	switch arg := vPtr.(type) {
	case (*Container):

		switch data[0] {
		case 2:
			type eType struct {
				Status string `json:"status"`
			}
			var e eType
			err := json.Unmarshal(data[1:], &e)
			arg.Err = errors.New(e.Status)
			return err

		default:
			i := n.mf.Build()
			err := json.Unmarshal(data, i)

			arg.I = i
			return err
		}
	case (*string):
		vPtr = string(data)
		return nil
	default:
		return n.je.Decode(subject, data, vPtr)
	}
}

// Container TODO
type Container struct {
	I   interface{}
	Err error
}

// RegisterEncoder ..
func RegisterEncoder(mf *dds.ModelFactory, entityName string) {
	nats.RegisterEncoder("dds."+entityName, &natsTypeEncoder{
		mf: mf,
	})
}

// Listen listens on nats for requests on the queue given the entity name
func Listen(ctx context.Context, mf *dds.ModelFactory, storage dds.Storage, entityName string) error {

	RegisterEncoder(mf, entityName)

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

	c, _ := nats.NewEncodedConn(nc, "dds."+entityName)

	go func() {
		<-ctx.Done()
		c.Close()
	}()

	c.Subscribe(entityName+".get.>", func(subj, reply string, data []byte) {

		items := strings.Split(subj, ".")
		id := items[len(items)-1]

		log.Printf("Getting %s %s", entityName, id)

		object, err := storage.Get(id)

		log.Printf("Results: %v %v", object, err)

		if err != nil {
			c.Publish(reply, statusResult{"status": err.Error()})
			return
		}

		if object == nil {
			c.Publish(reply, statusResult{"status": "NotFound"})
			return
		}

		c.Publish(reply, object)
	})

	c.Subscribe(entityName+".create", func(subj, reply string, i Container) {

		id, err := storage.Create(i.I)
		if err != nil {
			c.Publish(reply, map[string]string{
				"status": err.Error(),
			})
			return
		}

		c.Publish(reply, map[string]string{
			"id": id,
		})

		c.Publish("events."+entityName+".create", i.I)
	})

	c.Subscribe(entityName+".delete.>", func(subj, reply string, data []byte) {

		items := strings.Split(subj, ".")
		id := items[len(items)-1]

		err := storage.Delete(id)
		if err != nil {
			c.Publish(reply, map[string]string{"status": err.Error()})
			return
		}

		c.Publish(reply, map[string]string{"status": "OK"})

		c.Publish("events."+entityName+".delete", id)
	})

	c.Subscribe(entityName+".update", func(subj, reply string, i Container) {

		err := storage.Update(i.I)
		if err != nil {
			c.Publish(reply, statusResult{"status": err.Error()})
			return
		}

		c.Publish(reply, map[string]interface{}{"status": "OK"})

		c.Publish("events."+entityName+".update", i.I)
	})

	return nil
}

// Client creates a client
func Client(mf *dds.ModelFactory, entityName string) dds.Model {
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

	m.c.Subscribe(replyID, func(resp *Container) {
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

func (m *natsModel) Subscribe(evt string, out chan string) func() {
	sx, _ := m.c.Subscribe("events."+m.entityName+"."+evt, func(_ string, _ string, data []byte) {
		out <- string(data)
	})

	return func() {
		sx.Unsubscribe()
	}
}
