package natsSupport

import (
	"encoding/json"
	"strings"

	"github.com/CyCoreSystems/dds"
	"github.com/nats-io/nats"
)

func (m *natsModel) Subscribe(evt string) dds.Subscription {
	var nc natsSubscription
	nc.c = make(chan dds.Event, 10)
	sx, _ := m.c.Subscribe("events."+m.entityName+"."+evt, func(msg *nats.Msg) {

		subj := msg.Subject
		data := msg.Data

		items := strings.Split(subj, ".")

		nc.c <- dds.Event{
			Entity: m.entityName,
			Type:   items[2],
			Metadata: map[string]interface{}{
				"data": string(data),
			},
		}
	})

	nc.sx = sx

	return &nc
}

type natsSubscription struct {
	c  chan dds.Event
	sx *nats.Subscription
}

func (nc *natsSubscription) C() <-chan dds.Event {
	return nc.c
}

func (nc *natsSubscription) Close() {
	nc.sx.Unsubscribe()
	if nc.c != nil {
		close(nc.c)
	}
	nc.c = nil
}

func subscribe(conn *nats.Conn, errReporter func(err error, msg string), endpoint string, h Handler) (*nats.Subscription, error) {

	cb := func(msg *nats.Msg) {

		reply := msg.Reply
		data := msg.Data
		subj := msg.Subject

		h(subj, data, func(i interface{}, err error) {

			if err != nil {
				resp := []byte(err.Error())
				if err2 := conn.Publish(reply+".err", resp); err2 != nil {
					errReporter(err, "Error sending error reply")
				}
				return
			}

			resp := []byte("{}")
			if i != nil {
				resp, err = json.Marshal(i)
				if err != nil {
					errReporter(err, "Error building response reply")
					return
				}
			}

			if err = conn.Publish(reply+".resp", resp); err != nil {
				errReporter(err, "Error sending response reply")
			}
		})
	}

	sub, err := conn.Subscribe(endpoint, cb)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
