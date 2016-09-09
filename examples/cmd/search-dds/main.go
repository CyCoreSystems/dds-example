package main

import (
	"fmt"

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

	if err := natsSupport.ListenHandler(ctx, microblag.Search, searchHandler(0)); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
