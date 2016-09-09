package main

import (
	"fmt"

	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"

	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	storage := &entryStorage{
		data: make(map[string]microblag.Entry),
	}

	if err := natsSupport.Listen(ctx, microblag.EntryFactory, storage); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
