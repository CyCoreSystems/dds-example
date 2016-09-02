package nats

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/nats-io/nats"
)

// Listen listens on nats for requests on the queue given the entity name
func Listen(ctx context.Context, entityName string) error {

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

		//TODO: Read from MySQL
	})

	c.Subscribe(entityName+".create", func(subj, reply string, data string) {
		fmt.Printf("got create: %s %v\n", entityName, data)

		//TODO: Save to MySQL

		c.Publish(entityName+".events.create", "data")
	})

	c.Subscribe(entityName+".delete", func(subj, reply string, id string) {

		fmt.Printf("got delete: %s %v\n", entityName, id)

		//TODO: Delete from MySQL

		c.Publish(entityName+".events.delete", "data")
	})

	c.Subscribe(entityName+".update", func(subj, reply string, data string) {
		fmt.Printf("got update: %s %v\n", entityName, data)

		//TODO: Update in MySQL

		c.Publish(entityName+".events.update", data)
	})

	return nil
}
