package utils

import (
	"encoding/json"

	"github.com/rs/xid"
)

type XID string

func NewXID() XID {
	return XID(xid.New().Bytes())
}
func (id *XID) String() string {
	return string(*id)
}

func (id XID) MarshalBSON() ([]byte, error) {
	return json.Marshal(&id)
}

func (id *XID) UnmarshalBSON(data []byte) error {
	var tmp struct{ ID string }

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*id = XID(tmp.ID)

	return nil
}
