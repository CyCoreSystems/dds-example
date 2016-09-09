package natsSupport

import (
	"encoding/json"

	"github.com/nats-io/nats"
)

func publish(nc *nats.Conn, entityName string, event string, i interface{}) error {
	body, err := json.Marshal(i)
	if err != nil {
		return err
	}
	nc.Publish("events."+entityName+"."+event, body)
	return nil
}
