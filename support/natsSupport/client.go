package natsSupport

import (
	"errors"
	"fmt"
	"time"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
)

// Client creates a client
func Client(mf *dds.ModelFactory) dds.Model {

	entityName := mf.EntityName

	var nc *nats.Conn
	var err error

	for i := 0; i != 3 && nc == nil; i++ {
		<-time.After(500 * time.Millisecond)
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			fmt.Printf("Error connecting to nats: '%v'\n", err)
		}
	}

	if nc == nil {
		panic(errors.New("Failed to connect to nats, giving up\n"))
	}

	return &natsModel{
		c:          nc,
		entityName: entityName,
	}
}

type natsModel struct {
	c          *nats.Conn
	entityName string
}

func (m *natsModel) Close() {
	m.c.Close()
}

func (m *natsModel) Get(id string, i interface{}) (err error) {
	err = request(m.c, m.entityName+".get."+id, nil, &i)
	return
}

func (m *natsModel) Create(i interface{}) (str string, err error) {
	err = request(m.c, m.entityName+".create", i, &str)
	return
}

func (m *natsModel) Delete(id string) (err error) {
	err = request(m.c, m.entityName+".delete."+id, nil, nil)
	return
}

func (m *natsModel) Update(id string, i interface{}) (err error) {
	err = request(m.c, m.entityName+".update", i, nil)
	return
}
