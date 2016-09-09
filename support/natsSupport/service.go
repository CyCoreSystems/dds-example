package natsSupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
	"golang.org/x/net/context"
)

// Listen listens on nats for requests on the queue given the entity name
func Listen(ctx context.Context, mf *dds.ModelFactory, storage dds.Storage) error {

	entityName := mf.EntityName

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

	go func() {
		<-ctx.Done()
		nc.Close()
	}()

	var reporter = func(err error, msg string) {
		fmt.Printf("ERR: %s: '%v'", msg, err)
	}

	subscribe(nc, reporter, entityName+".get.>", func(subj string, request []byte, reply Reply) {

		items := strings.Split(subj, ".")
		id := items[len(items)-1]

		log.Printf("Getting %s %s", entityName, id)

		object, err := storage.Get(id)

		log.Printf("Results: %v %v", object, err)

		if err != nil {
			reply(nil, err)
			return
		}

		if object == nil {
			reply(nil, errors.New("Not Found"))
			return
		}

		reply(object, nil)
	})

	subscribe(nc, reporter, entityName+".create", func(subj string, request []byte, reply Reply) {
		i := mf.Build()
		id, err := storage.Create(i)
		reply(id, err)

		if err == nil {
			publish(nc, entityName, "create", i)
		}
	})

	subscribe(nc, reporter, entityName+".delete.>", func(subj string, _ []byte, reply Reply) {

		items := strings.Split(subj, ".")
		id := items[len(items)-1]

		err := storage.Delete(id)
		reply(nil, err)

		if err == nil {
			publish(nc, entityName, "delete", id)
		}
	})

	subscribe(nc, reporter, entityName+".update", func(subj string, body []byte, reply Reply) {
		i := mf.Build()
		err := json.Unmarshal(body, i)
		if err != nil {
			reply(nil, err)
			return
		}

		err = storage.Update(i)
		reply(nil, err)

		if err == nil {
			publish(nc, entityName, "update", i)
		}
	})

	return nil
}
