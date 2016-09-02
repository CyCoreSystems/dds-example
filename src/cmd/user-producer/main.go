package main

import (
	"errors"
	"fmt"
	"microblag"
	"time"

	"golang.org/x/net/context"

	"github.com/nats-io/nats"
	"github.com/satori/go.uuid"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user microblag.User
	user.Username = "X"
	user.Created = time.Now()

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

	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	replyID := uuid.NewV1().String()

	c.Subscribe(replyID, func(v interface{}) {
		fmt.Printf("got create response: %v\n", v)
		cancel()
	})

	c.PublishRequest("users.create", replyID, &user)

	<-ctx.Done()

	fmt.Printf("err: %v\n", ctx.Err())
}
