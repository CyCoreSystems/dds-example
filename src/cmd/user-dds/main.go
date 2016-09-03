package main

import (
	"dds"
	"dds/nats"
	"fmt"
	"microblag"

	"golang.org/x/net/context"
)

var users dds.Model

func main() {

	ctx := context.Background()

	storage := &userStorage{
		data: make(map[string]microblag.User),
	}

	if err := dnats.Listen(ctx, microblag.UserFactory, storage); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
