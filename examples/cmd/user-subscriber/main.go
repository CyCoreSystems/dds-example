package main

import (
	"fmt"

	"github.com/CyCoreSystems/dds"
	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"
)

func main() {

	client := natsSupport.Client(microblag.UserFactory)
	defer client.Close()

	s1 := client.Subscribe(dds.CreateEvent)
	defer s1.Close()

	s2 := client.Subscribe(dds.AllEvents)
	defer s2.Close()

	for {
		select {
		case buf := <-s1.C():
			fmt.Printf("got create event: %v\n", buf)
		case buf := <-s2.C():
			fmt.Printf("got all event: %v\n", buf)
		}
	}

}
