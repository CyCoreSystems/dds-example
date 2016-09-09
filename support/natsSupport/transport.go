package natsSupport

import (
	"errors"
	"fmt"
	"time"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
)

type natsTransport struct {
	nc *nats.Conn
}

func (nt *natsTransport) Close() {
	nt.nc.Close()
}

// NewTransport builds a new NATS transport
func NewTransport() (dds.Transport, error) {
	var nc *nats.Conn
	var err error

	for i := 0; i != 3 && nc == nil; i++ {
		<-time.After(500 * time.Millisecond)
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			fmt.Printf("Error connecting to nats: '%v'", err)
		}
	}

	if nc == nil {
		return nil, errors.New("Failed to connect to nats, giving up\n")
	}

	return &natsTransport{
		nc: nc,
	}, nil
}
