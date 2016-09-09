package main

import (
	"fmt"

	"github.com/CyCoreSystems/dds"
	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"

	"golang.org/x/net/context"
)

type searchHandler int

func (sh searchHandler) Handle(request interface{}) (response interface{}, err error) {
	s := request.(*microblag.SearchRequest)
	fmt.Printf("got search: %v", s)

	return []string{"hello", "world"}, nil
}

func main() {

	ctx := context.Background()

	nt, err := natsSupport.NewTransport()
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	svc := dds.NewActionService(microblag.Search, searchHandler(0), nt)

	if err := svc.Listen(); err != nil {
		fmt.Printf("Err: %v\n", err)
	}

	<-ctx.Done()
}
