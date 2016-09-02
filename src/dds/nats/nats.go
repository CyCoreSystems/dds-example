package nats

import (
	"dds"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/nats-io/nats"
)

// Listen listens on nats for requests on the queue given the entity name
func Listen(ctx context.Context, storage dds.Storage, entityName string) error {

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
		c.Close()
	}()

	c.Subscribe(entityName+".get.>", func(subj, reply string) {
		fmt.Printf("got read: %s %v\n", entityName, subj)

		items := strings.Split(subj, ".")
		id := items[len(items)-1]

		object, err := storage.Get(id)
		if err != nil {
			c.Publish(reply, err)
			return
		}

		c.Publish(reply, object)
	})

	c.Subscribe(entityName+".create", func(subj, reply string, i interface{}) {
		fmt.Printf("got create: %s %v\n", entityName, i)

		id, err := storage.Create(i)
		if err != nil {
			c.Publish(reply, map[string]string{
				"err": err.Error(),
			})
			return
		}

		c.Publish(reply, map[string]string{
			"id": id,
		})

		c.Publish("events."+entityName+".create", i)
	})

	c.Subscribe(entityName+".delete", func(subj, reply string, id string) {

		fmt.Printf("got delete: %s %v\n", entityName, id)

		err := storage.Delete(id)
		if err != nil {
			c.Publish(reply, err)
			return
		}

		c.Publish(reply, "OK")

		c.Publish("events."+entityName+".delete", id)
	})

	c.Subscribe(entityName+".update", func(subj, reply string, i interface{}) {
		fmt.Printf("got update: %s %v\n", entityName, i)

		err := storage.Update(i)
		if err != nil {
			c.Publish(reply, err)
			return
		}

		c.Publish(reply, "OK")

		c.Publish("events."+entityName+".update", i)
	})

	return nil
}
