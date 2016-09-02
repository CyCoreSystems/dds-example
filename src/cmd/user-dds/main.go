package main

import (
	"dds/nats"
	"fmt"
	"microblag"

	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	storage := &userStorage{
		data: make(map[string]microblag.User),
	}

	if err := nats.Listen(ctx, storage, "users"); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
