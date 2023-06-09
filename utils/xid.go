package utils

import (
	"encoding/json"

	"github.com/rs/xid"
)

type XID string

func NewXID() XID {
	return XID(xid.New().String())
}
func (id *XID) String() string {
	return string(*id)
}

func (id XID) MarshalBSON() ([]byte, error) {
	return json.Marshal(id)
}

func (id *XID) UnmarshalBSON(data []byte) error {
	var xid string

	if err := json.Unmarshal(data, &xid); err != nil {
		return err
	}

	*id = XID(xid)

	return nil
}
