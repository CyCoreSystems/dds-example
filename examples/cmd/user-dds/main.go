package main

import (
	"fmt"

	"github.com/CyCoreSystems/dds"
	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"

	"golang.org/x/net/context"
)

var users dds.Model

func main() {

	ctx := context.Background()

	storage := &userStorage{
		data: make(map[string]microblag.User),
	}

	if err := natsSupport.Listen(ctx, microblag.UserFactory, storage); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
