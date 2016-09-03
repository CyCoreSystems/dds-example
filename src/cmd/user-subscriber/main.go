package main

import (
	"fmt"
	"microblag"

	dnats "dds/nats"
)

func main() {

	client := dnats.Client(microblag.UserFactory, "users")
	defer client.Close()

	data := make(chan string)
	closer := client.Subscribe("create", data)
	defer closer()

	closer2 := client.Subscribe("*", data)
	defer closer2()

	for {
		select {
		case buf := <-data:
			fmt.Printf("got event: %s\n", buf)
		}
	}

}
