package main

import (
	"dds/nats"
	"fmt"

	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	if err := nats.Listen(ctx, "entries"); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
