package main

import (
	"dds/nats"
	"fmt"
	"microblag"

	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	storage := &entryStorage{
		data: make(map[string]microblag.Entry),
	}

	if err := nats.Listen(ctx, microblag.EntryFactory, storage, "entries"); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
