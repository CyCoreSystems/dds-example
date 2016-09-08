package main

import (
	"dds/nats"
	"fmt"
	"microblag"

	"golang.org/x/net/context"
)

type searchHandler int

func (_ searchHandler) Handle(request interface{}) (response interface{}, err error) {
	s := request.(*microblag.SearchRequest)
	fmt.Printf("got search: %v", s)

	return []string{"hello", "world"}, nil
}

func main() {

	ctx := context.Background()

	if err := dnats.ListenHandler(ctx, microblag.Search, searchHandler(0)); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	<-ctx.Done()
}
