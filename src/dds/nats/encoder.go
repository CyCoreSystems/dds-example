package dnats

import (
	"dds"
	"encoding/json"
	"errors"

	"github.com/nats-io/nats"
	"github.com/nats-io/nats/encoders/builtin"
)

// RegisterEncoder ..
func RegisterEncoder(mf *dds.ModelFactory, entityName string) {
	nats.RegisterEncoder("dds."+entityName, &natsTypeEncoder{
		mf: mf,
	})
}

type statusResult map[string]interface{}

type natsTypeEncoder struct {
	mf *dds.ModelFactory
	je builtin.JsonEncoder
}

func (n *natsTypeEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	switch v.(type) {
	case statusResult:
		b, err := n.je.Encode(subject, v)
		b = append([]byte{2}, b...)
		return b, err
	default:
		return n.je.Encode(subject, v)
	}
}

func (n *natsTypeEncoder) Decode(subject string, data []byte, vPtr interface{}) error {

	switch arg := vPtr.(type) {
	case (*container):

		switch data[0] {
		case 2:
			type eType struct {
				Status string `json:"status"`
			}
			var e eType
			err := json.Unmarshal(data[1:], &e)
			arg.Err = errors.New(e.Status)
			return err

		default:
			i := n.mf.Build()
			err := json.Unmarshal(data, i)

			arg.I = i
			return err
		}
	case (*string):
		vPtr = string(data)
		return nil
	default:
		return n.je.Decode(subject, data, vPtr)
	}
}
