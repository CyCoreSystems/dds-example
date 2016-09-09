package main

import (
	"fmt"
	"time"

	"github.com/CyCoreSystems/dds/examples/microblag"
	"github.com/CyCoreSystems/dds/support/natsSupport"
)

func main() {

	var user microblag.User
	user.Username = "X"
	user.Created = time.Now()

	client := natsSupport.Client(microblag.UserFactory)
	defer client.Close()

	search := natsSupport.ActionClient(microblag.Search)

	resp, err := search(&microblag.SearchRequest{})
	if err != nil {
		fmt.Printf("Error searching: %s\n", err)
		return
	}

	fmt.Printf("search resp: %v\n", resp)

	id, err := client.Create(user)
	if err != nil {
		fmt.Printf("Error creating user: %s\n", err)
		return
	}

	fmt.Printf("Got ID: %s\n", id)

	var myUser *microblag.User
	err = client.Get(id, &myUser)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
	}

	fmt.Printf("user : %v\n", myUser)

	var d = "hello-world"
	myUser.DisplayName = &d
	err = client.Update(id, myUser)
	if err != nil {
		fmt.Printf("Error updating user: %s\n", err)

	}

	fmt.Printf("user : %v\n", myUser)

	err = client.Get(id, &myUser)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
	}

	fmt.Printf("user : %v\n", myUser)

	err = client.Delete(id)
	if err != nil {
		fmt.Printf("Error deleting user: %s\n", err)
	}

	err = client.Get(id, &myUser)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
	}

}
