package dnats

import (
	"dds"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats"
	"golang.org/x/net/context"
)

// Listen listens on nats for requests on the queue given the entity name
func Listen(ctx context.Context, mf *dds.ModelFactory, storage dds.Storage) error {

	entityName := mf.EntityName

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

	c.Subscribe(entityName+".create", func(subj, reply string, i container) {

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

	c.Subscribe(entityName+".update", func(subj, reply string, i container) {

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
