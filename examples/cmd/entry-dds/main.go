package main

import (
	"fmt"

	"github.com/CyCoreSystems/dds"
	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"

	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	storage := &entryStorage{
		data: make(map[string]microblag.Entry),
	}

	transport, err := natsSupport.NewTransport()
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	svc := dds.NewDataService(microblag.EntryFactory, storage, transport)
	defer svc.Stop()

	err = svc.Listen()
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
